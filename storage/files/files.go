package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"memolink-bot/storage"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774 //read and write

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return fmt.Errorf("cant save: %w", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("cant get file name: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("cant create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("cant encode data: %w", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cant read from directory: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("cant get fileName: %w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("cant remove file: %w", err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("cant get fileName: %w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("cant check is file exists: %w", err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cant open file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("cant decode page: %w", err)
	}

	return &p, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
