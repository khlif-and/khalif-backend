package storage

import (
	"fmt"
	"io"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct{}

func NewLocalStorage() ports.StorageProvider {
	utils.InitUploadDirs()
	return &LocalStorage{}
}

func (s *LocalStorage) UploadFile(file io.Reader, filename string, subdirectory string) (string, string, error) {
	// Generate Unique Filename
	newFilename := utils.GenerateUniqueFilename(filename)
	destDir := filepath.Join(utils.UploadDir, subdirectory)
	destPath := filepath.Join(destDir, newFilename)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Stream copy
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(destPath)
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	// Calculate relative URL
	relativeURL := fmt.Sprintf("/%s/%s/%s", strings.TrimPrefix(utils.UploadDir, "./"), subdirectory, newFilename)

	return relativeURL, destPath, nil
}

func (s *LocalStorage) DeleteFile(path string) error {
	return utils.DeleteUploadedFile(path)
}
