package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStorage struct {
	Dir string
}

func NewLocalStorage(dir string) *LocalStorage {
	return &LocalStorage{Dir: dir}
}

// Save saves an uploaded file stream; returns filename saved (basename)
func (s *LocalStorage) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	// extension (keep the original extension)
	ext := filepath.Ext(header.Filename)
	id := uuid.New().String()
	name := id + ext
	destinationPath := filepath.Join(s.Dir, name)

	out, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open destination: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}

	return name, nil
}
