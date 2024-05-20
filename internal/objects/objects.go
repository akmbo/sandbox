package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
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

type blobObject struct {
	DateType string
	Size     int
	Reader   io.Reader
}

func getObjectReader(repo repository.Repository, checksum string) (*os.File, error) {
	if len(checksum) != 40 {
		return nil, errors.New("string provided is not a valid checksum")
	}

	obj_path := filepath.Join(repo.Objects, checksum[:2], checksum[2:])
	_, err := os.Stat(obj_path)
	if err != nil {
		return nil, err
	}

	r, err := os.Open(obj_path)
	if err != nil {
		return nil, err
	}

	return r, err
}

func ReadBlobV2(repo repository.Repository, checksum string) (blob blobObject, err error) {
	file, err := getObjectReader(repo, checksum)
	if err != nil {
		return blob, err
	}
	defer file.Close()
	zr, err := zlib.NewReader(file)
	if err != nil {
		return blob, err
	}
	defer zr.Close()

	rawHeader := make([]byte, 1024)
	_, err = zr.Read(rawHeader)
	if err != nil {
		return blob, err
	}

	getObjectType := func(b []byte) (objectType string, remain []byte) {
		var objectTypeEnd int
		for i, r := range b {
			if r == ' ' {
				objectTypeEnd = i
				break
			}
		}
		return fmt.Sprint(b[:objectTypeEnd]), b[:objectTypeEnd]
	}

	getContentSize := func(b []byte) (contentSize int, remain []byte, err error) {
		var contentSizeEnd int
		for i, r := range b {
			if r == '\u0000' {
				contentSizeEnd = i
				break
			}
		}
		sizeString := fmt.Sprint(b[:contentSizeEnd])
		size, err := strconv.Atoi(sizeString)
		if err != nil {
			return contentSize, remain, err
		}
		return size, b[contentSizeEnd:], nil
	}

	objecType, remain := getObjectType(rawHeader)
	contentSize, _, err := getContentSize(remain)
	if err != nil {
		return blob, err
	}

	// TODO: need to combine the remaining bytes and the remaining zlib reader into one reader
	blob = blobObject{objecType, contentSize, zr}

	return blob, nil
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

	obj_path := filepath.Join(repo.Objects, checksum[:2], checksum[2:])
	_, err = os.Stat(obj_path)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(obj_path)
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
