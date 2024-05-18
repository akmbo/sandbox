package repository

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type repository struct {
	Path    string
	Data    string
	Config  string
	Head    string
	Hooks   string
	Info    string
	Objects string
	Refs    string
}

func newRepository(projectPath string) repository {
	repoDir := filepath.Join(projectPath, ".minigit")
	r := repository{
		Path:    projectPath,
		Data:    repoDir,
		Config:  filepath.Join(repoDir, "config"),
		Head:    filepath.Join(repoDir, "HEAD"),
		Hooks:   filepath.Join(repoDir, "hooks"),
		Info:    filepath.Join(repoDir, "info"),
		Objects: filepath.Join(repoDir, "objects"),
		Refs:    filepath.Join(repoDir, "refs"),
	}
	return r
}

func walkUpSearch(path string, dir string) (string, error) {
	_, err := os.Stat(filepath.Join(path, dir))
	if errors.Is(err, fs.ErrNotExist) {
		if path == "/" {
			return "", err
		}
		return walkUpSearch(filepath.Dir(path), dir)
	}
	if err != nil {
		return "", err
	}
	return path, nil
}

func Discover(path string) (repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return repository{}, err
	}

	if _, err = os.Stat(absPath); err != nil {
		return repository{}, err
	}

	projectPath, err := walkUpSearch(path, ".minigit")
	if err != nil {
		return repository{}, err
	}

	r := newRepository(projectPath)

	return r, nil
}

func Create(path string) (repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return repository{}, err
	}

	if _, err = os.Stat(absPath); err != nil {
		return repository{}, err
	}

	r := newRepository(absPath)

	makeDirs := func(dirs ...string) error {
		for _, d := range dirs {
			err := os.MkdirAll(d, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	}

	makeFiles := func(paths ...string) error {
		for _, p := range paths {
			f, err := os.Create(p)
			if err != nil {
				return err
			}
			f.Close()
		}
		return nil
	}

	err = makeDirs(r.Data, r.Hooks, r.Info, r.Objects, r.Refs)
	if err != nil {
		return repository{}, err
	}
	err = makeFiles(r.Config, r.Head)
	if err != nil {
		return repository{}, err
	}

	return r, nil
}
