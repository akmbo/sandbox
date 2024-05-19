package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aaolen/mini-git/internal/repository"
)

func compressString(content string) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(content))
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompressString(compressed []byte) (string, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	if err := r.Close(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func getChecksum(content string) string {
	h := sha1.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func WriteBlob(repo repository.Repository, content string) (checksum string, err error) {
	compressed, err := compressString(content)
	if err != nil {
		return "", err
	}
	checksum = getChecksum(content)

	err = os.MkdirAll(filepath.Join(repo.Objects, checksum[:2]), 0777)
	if err != nil {
		return "", err
	}

	objPath := filepath.Join(repo.Objects, checksum[:2], checksum[2:])

	err = os.WriteFile(objPath, compressed, 0644)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

func ReadBlob(repo repository.Repository, checksum string) (content string, err error) {
	obj_path := filepath.Join(repo.Objects, checksum[:2], checksum[2:])
	_, err = os.Stat(obj_path)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(obj_path)
	if err != nil {
		return "", err
	}

	content, err = decompressString(b)
	if err != nil {
		return "", err
	}

	return content, nil
}
