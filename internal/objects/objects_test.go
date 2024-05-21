package objects

import (
	"io"
	"os"
	"testing"

	"github.com/aaolen/mini-git/internal/repository"
	"github.com/stretchr/testify/assert"
)

func makeTempRepo() (repo repository.Repository, dir string, err error) {
	dir, err = os.MkdirTemp("", "*")
	if err != nil {
		return repo, dir, err
	}

	repo, err = repository.Create(dir)
	if err != nil {
		return repo, dir, err
	}

	return repo, dir, err
}

func TestReadAndWriteBlob(t *testing.T) {
	assert := assert.New(t)

	r, dir, err := makeTempRepo()
	assert.NoError(err)
	defer os.RemoveAll(dir)

	inputContent := "hello world"

	checksum, err := WriteBlob(r, inputContent)
	assert.NoError(err)

	outputContent, err := ReadBlob(r, checksum)
	assert.NoError(err)

	assert.Equal(inputContent, outputContent)
}

func TestReadHeader(t *testing.T) {
	assert := assert.New(t)

	r, dir, err := makeTempRepo()
	defer os.RemoveAll(dir)
	assert.NoError(err)

	checksum, err := WriteBlob(r, "hello world")
	assert.NoError(err)

	expected := Header{"blob", 11}

	result, err := ReadHeader(r, checksum)
	assert.NoError(err)

	assert.Equal(expected, result)
}

func TestGetContentReader(t *testing.T) {
	assert := assert.New(t)

	r, dir, err := makeTempRepo()
	assert.NoError(err)
	defer os.RemoveAll(dir)

	checksum, err := WriteBlob(r, "hello world")
	assert.NoError(err)

	reader, close, err := GetContentReader(r, checksum)
	defer close()
	assert.NoError(err)

	resultBytes, err := io.ReadAll(reader)
	assert.NoError(err)

	expected := "hello world"
	result := string(resultBytes)

	assert.Equal(expected, result)
}
