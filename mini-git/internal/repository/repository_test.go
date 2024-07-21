package repository

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeDirs(dirs ...string) error {
	for _, d := range dirs {
		err := os.MkdirAll(d, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test_walkUpSearch(t *testing.T) {
	assert := assert.New(t)

	t.Run("Directory with desired path inside", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "*")
		assert.NoError(err)
		defer os.RemoveAll(dir)

		searchedForDir := filepath.Join(dir, "find_me")

		err = makeDirs(searchedForDir)
		assert.NoError(err)

		result, err := walkUpSearch(dir, "find_me")
		assert.NoError(err)
		assert.Equal(dir, result)
	})

	t.Run("Directory several layers deeper than desired path", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "*")
		assert.NoError(err)

		searchedForDir := filepath.Join(dir, "find_me")
		nestedPath := filepath.Join(dir, "one", "two", "three")

		err = makeDirs(searchedForDir, nestedPath)
		assert.NoError(err)

		result, err := walkUpSearch(nestedPath, "find_me")
		assert.NoError(err)
		assert.Equal(dir, result)
	})

	t.Run("No parent path containing desired path", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "*")
		assert.NoError(err)
		defer os.RemoveAll(dir)

		nestedPath := filepath.Join(dir, "one", "two", "three")

		err = makeDirs(nestedPath)
		assert.NoError(err)

		_, err = walkUpSearch(nestedPath, "find_me")
		assert.True(errors.Is(err, fs.ErrNotExist))
	})

	t.Run("Random non-existent path", func(t *testing.T) {
		randomPath := "/tmp/this/path/probably/does/not/exist/12345"
		_, err := walkUpSearch(randomPath, "find_me")
		assert.True(errors.Is(err, fs.ErrNotExist))
	})

}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "*")
	assert.NoError(err)
	defer os.RemoveAll(dir)

	r, err := Create(dir)
	assert.NoError(err)

	validatePaths := func(paths ...string) error {
		for _, p := range paths {
			_, err = os.Stat(p)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = validatePaths(
		r.Path,
		r.Data,
		r.Config,
		r.Head,
		r.Hooks,
		r.Info,
		r.Objects,
		r.Refs,
	)
	assert.NoError(err)
}
