package objects

import (
	"os"
	"testing"

	"github.com/aaolen/mini-git/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestReadAndWriteBlob(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "*")
	assert.NoError(err)
	defer os.RemoveAll(dir)

	r, err := repository.Create(dir)
	assert.NoError(err)

	inputContent := "hello world"

	checksum, err := WriteBlob(r, inputContent)
	assert.NoError(err)

	outputContent, err := ReadBlob(r, checksum)
	assert.NoError(err)

	assert.Equal(inputContent, outputContent)
}
