package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxFileSize    = 5 * 1024 * 1024 // 5MB
	MaxAudioSize   = 50 * 1024 * 1024 // 50MB for audio files
	UploadDir      = "./uploads"
	ProfilePicDir  = "profile_pictures"
	AudioDir       = "audio"
	ThumbnailDir   = "thumbnails"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type UploadResult struct {
	Filename string
	Path     string
	URL      string
	Size     int64
}

func InitUploadDirs() error {
	dirs := []string{
		filepath.Join(UploadDir, ProfilePicDir),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory %s: %w", dir, err)
		}
	}
	return nil
}

func ValidateImageFile(file *multipart.FileHeader) error {
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed (%d MB)", MaxFileSize/(1024*1024))
	}

	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return fmt.Errorf("invalid file type: %s. Allowed: jpeg, png, gif, webp", contentType)
	}

	return nil
}

func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%d_%s%s", timestamp, uniqueID, ext)
}

func SaveUploadedFile(file *multipart.FileHeader, subDir string) (*UploadResult, error) {
	if err := ValidateImageFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	filename := GenerateUniqueFilename(file.Filename)
	destDir := filepath.Join(UploadDir, subDir)
	destPath := filepath.Join(destDir, filename)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	relativeURL := fmt.Sprintf("/%s/%s/%s", strings.TrimPrefix(UploadDir, "./"), subDir, filename)

	return &UploadResult{
		Filename: filename,
		Path:     destPath,
		URL:      relativeURL,
		Size:     written,
	}, nil
}

func DeleteUploadedFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	
	return os.Remove(filePath)
}

func GetProfilePicturePath(filename string) string {
	return filepath.Join(UploadDir, ProfilePicDir, filename)
}

var AllowedAudioTypes = map[string]bool{
	"audio/mpeg":  true, // mp3
	"audio/mp3":   true,
	"audio/wav":   true,
	"audio/x-wav": true,
	"audio/ogg":   true,
	"audio/aac":   true,
	"audio/m4a":   true,
	"audio/x-m4a": true,
}

func ValidateAudioFile(file *multipart.FileHeader) error {
	if file.Size > MaxAudioSize {
		return fmt.Errorf("audio file size exceeds maximum allowed (%d MB)", MaxAudioSize/(1024*1024))
	}

	contentType := file.Header.Get("Content-Type")
	if !AllowedAudioTypes[contentType] {
		return fmt.Errorf("invalid audio file type: %s. Allowed: mp3, wav, ogg, aac, m4a", contentType)
	}

	return nil
}

func SaveAudioFile(file *multipart.FileHeader) (*UploadResult, error) {
	if err := ValidateAudioFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	filename := GenerateUniqueFilename(file.Filename)
	destDir := filepath.Join(UploadDir, AudioDir)
	destPath := filepath.Join(destDir, filename)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	relativeURL := fmt.Sprintf("/%s/%s/%s", strings.TrimPrefix(UploadDir, "./"), AudioDir, filename)

	return &UploadResult{
		Filename: filename,
		Path:     destPath,
		URL:      relativeURL,
		Size:     written,
	}, nil
}
