package ports

import (
	"io"
)

type StorageProvider interface {
	UploadFile(file io.Reader, filename string, directory string) (string, string, error)
	DeleteFile(path string) error
}
