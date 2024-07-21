package repository

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Repository struct {
	Path    string
	Data    string
	Config  string
	Head    string
	Hooks   string
	Info    string
	Objects string
	Refs    string
}

func newRepository(projectPath string) Repository {
	repoDir := filepath.Join(projectPath, ".minigit")
	r := Repository{
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

func Discover(path string) (Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Repository{}, err
	}

	if _, err = os.Stat(absPath); err != nil {
		return Repository{}, err
	}

	projectPath, err := walkUpSearch(absPath, ".minigit")
	if err != nil {
		return Repository{}, err
	}

	r := newRepository(projectPath)

	return r, nil
}

func Create(path string) (Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Repository{}, err
	}

	if _, err = os.Stat(absPath); err != nil {
		return Repository{}, err
	}

	if _, err = os.Stat(filepath.Join(absPath, ".minigit")); err == nil {
		return Repository{}, errors.New("repository already initialized in directory")
	}

	r := newRepository(absPath)

	makeDirs := func(dirs ...string) error {
		for _, d := range dirs {
			err := os.MkdirAll(d, 0777)
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
		return Repository{}, err
	}
	err = makeFiles(r.Config, r.Head)
	if err != nil {
		return Repository{}, err
	}

	return r, nil
}
