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

	"github.com/aaolen/mini-git/internal/repository"
)

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
