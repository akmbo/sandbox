package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aaolen/mini-git/internal/repository"
)

// READ OBJECT
// take hash, open file reader
// transform into zlib writer
// create object struct
// read to first space, add type to object struct
// read to null byte, add size to object struct
// add reader with remaining bytes to struct

// WRITE OBJECT
// take reader, transform into multi-reader for hash, zlib, and size
// get size
// write type and size to header
// create zlib writer with header and file
// get hash
// write to file with hash name and compressed content

type Header struct {
	DataType string
	Size     int64
}

func getObjectPath(repo repository.Repository, checksum string) (string, error) {
	if len(checksum) != 40 {
		return "", errors.New("string provided is not a valid checksum")
	}

	objPath := filepath.Join(repo.Objects, checksum[:2], checksum[2:])
	return objPath, nil
}

func buildCompressedFileReader(path string) (r io.ReadCloser, closer func(), err error) {
	r, err = os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	zr, err := zlib.NewReader(r)
	if err != nil {
		r.Close()
		return nil, nil, err
	}

	closer = func() {
		zr.Close()
		r.Close()
	}

	return zr, closer, err
}

func buildHeader(b []byte) (Header, error) {
	var objectTypeEnd int
	var contentSizeEnd int

	s := string(b)

	for i, c := range s {
		if objectTypeEnd == 0 {
			if c == ' ' {
				objectTypeEnd = i
			}
			continue
		}
		if contentSizeEnd == 0 {
			if c == '\u0000' {
				contentSizeEnd = i
				break
			}
		}
	}

	objectType := s[:objectTypeEnd]
	contentSizeString := s[objectTypeEnd+1 : contentSizeEnd]

	contentSize, err := strconv.ParseInt(contentSizeString, 10, 64)
	if err != nil {
		return Header{}, err
	}

	h := Header{objectType, contentSize}

	return h, nil
}

func ReadHeader(repo repository.Repository, checksum string) (Header, error) {
	objPath, err := getObjectPath(repo, checksum)
	if err != nil {
		return Header{}, err
	}
	zr, close, err := buildCompressedFileReader(objPath)
	if err != nil {
		return Header{}, err
	}
	defer close()

	containingHeader := make([]byte, 64)
	n, err := io.ReadFull(zr, containingHeader)
	if err == io.EOF && n > 0 {
		containingHeader = containingHeader[:n]
	} else if err == io.ErrUnexpectedEOF {
		containingHeader = containingHeader[:n]
	} else if err != nil {
		return Header{}, err
	}

	h, err := buildHeader(containingHeader)
	if err != nil {
		return Header{}, err
	}

	return h, nil
}

func removeHeaderV2(b []byte) ([]byte, error) {
	var contentStart int
	for i, c := range b {
		if c == '\u0000' {
			contentStart = i + 1
		}
	}
	if contentStart == 0 {
		return nil, errors.New("Header end not found")
	}
	withoutHeader := b[contentStart:]
	return withoutHeader, nil
}

func GetContentReader(repo repository.Repository, checksum string) (reader io.Reader, closer func(), err error) {
	objPath, err := getObjectPath(repo, checksum)
	if err != nil {
		return reader, closer, err
	}
	zr, close, err := buildCompressedFileReader(objPath)
	if err != nil {
		return reader, closer, err
	}

	containingHeader := make([]byte, 64)
	n, err := io.ReadFull(zr, containingHeader)
	if err == io.EOF && n > 0 {
		containingHeader = containingHeader[:n]
	} else if err == io.ErrUnexpectedEOF {
		containingHeader = containingHeader[:n]
	} else if err != nil {
		return reader, closer, err
	}

	remaining, err := removeHeaderV2(containingHeader)
	if err != nil {
		return reader, closer, err
	}

	br := bytes.NewReader(remaining)
	mr := io.MultiReader(br, zr)

	return mr, close, err
}

// hash the raw contents from the given reader
// create the header, create a multi-reader with the header and the contentReader
// create a zlib writer that takes the multi-reader
// open a new file for writing with the hash
// write to the file with the zlib writer

func hashAndReuse(r io.Reader, w io.Writer) (string, error) {
	hasher := sha1.New()
	teeReader := io.TeeReader(r, hasher)
	if _, err := io.Copy(w, teeReader); err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	hashString := hex.EncodeToString(hash)
	return hashString, nil
}

func prefixHeader(r io.Reader) (io.Reader, error) {
	// var b bytes.Buffer
	// b.WriteString("blob ")
	// b.WriteString(fmt.Sprint(len(content)))
	// b.WriteString("\u0000")
	// b.Write(content)

	mr := io.MultiReader(r)
	return mr, nil
}

func buildCompressedFileWriter(path string) (w io.Writer, closer func(), err error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	zw := zlib.NewWriter(file)

	closer = func() {
		zw.Close()
		file.Close()
	}

	return zw, closer, err
}

func WriteContentWithHeader(repo repository.Repository, contentReader io.Reader) (checksum string, err error) {
	var content bytes.Buffer
	checksum, err = hashAndReuse(contentReader, &content)
	if err != nil {
		return checksum, err
	}

	withHeader, err := prefixHeader(&content)
	if err != nil {
		return checksum, err
	}

	objPath, err := getObjectPath(repo, checksum)
	if err != nil {
		return checksum, err
	}

	outFile, close, err := buildCompressedFileWriter(objPath)
	defer close()
	if err != nil {
		return checksum, err
	}

	buffer := make([]byte, 4096)

	for {
		n, err := withHeader.Read(buffer)
		if err != nil && err != io.EOF {
			return checksum, err
		}
		if n == 0 {
			break
		}
		_, err = outFile.Write(buffer[:n])
		if err != nil {
			return checksum, err
		}
	}

	return checksum, nil
}

func compress(input []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompress(compressed []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	if err := r.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func getChecksum(content string) string {
	h := sha1.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func addHeader(content []byte) []byte {
	var b bytes.Buffer
	b.WriteString("blob ")
	b.WriteString(fmt.Sprint(len(content)))
	b.WriteString("\u0000")
	b.Write(content)
	return b.Bytes()
}

func removeHeader(content []byte) []byte {
	var contentStart int
	for i, c := range content {
		if c == '\u0000' {
			contentStart = i + 1
			break
		}
	}
	return content[contentStart:]
}

func WriteBlob(repo repository.Repository, content string) (checksum string, err error) {
	checksum = getChecksum(content)
	withHeader := addHeader([]byte(content))
	store, err := compress(withHeader)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(repo.Objects, checksum[:2]), 0777)
	if err != nil {
		return "", err
	}

	objPath := filepath.Join(repo.Objects, checksum[:2], checksum[2:])

	err = os.WriteFile(objPath, store, 0644)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

func ReadBlob(repo repository.Repository, checksum string) (content string, err error) {
	if len(checksum) != 40 {
		return "", errors.New("string provided is not a valid checksum")
	}

	objPath := filepath.Join(repo.Objects, checksum[:2], checksum[2:])
	_, err = os.Stat(objPath)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(objPath)
	if err != nil {
		return "", err
	}

	output, err := decompress(b)
	if err != nil {
		return "", err
	}

	withoutHeader := removeHeader(output)

	return string(withoutHeader), nil
}
