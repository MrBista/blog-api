package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type StorageService interface {
	SaveFile(file *multipart.FileHeader, subPath string) (string, error)
	DeleteFile(uri string) error
}

type LocalStorage struct {
	BasePath string
	BaseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	_ = os.MkdirAll(basePath, os.ModePerm)
	return &LocalStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

func saveMultipartFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *LocalStorage) SaveFile(file *multipart.FileHeader, subPath string) (string, error) {

	subPath = strings.Trim(subPath, "/")

	targetDir := filepath.Join(s.BasePath, filepath.FromSlash(subPath))
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join(targetDir, filename)

	if err := saveMultipartFile(file, savePath); err != nil {
		return "", err
	}

	base := strings.TrimRight(s.BaseURL, "/")
	parts := []string{base}
	if subPath != "" {
		parts = append(parts, path.Clean(subPath))
	}
	parts = append(parts, url.PathEscape(filename))
	fileURL := strings.Join(parts, "/")

	return fileURL, nil
}

func (s *LocalStorage) DeleteFile(uri string) error {
	base := strings.TrimRight(s.BaseURL, "/")

	if strings.HasPrefix(uri, base+"/") || uri == base {
		rel := strings.TrimPrefix(uri, base)
		rel = strings.TrimLeft(rel, "/")

		// decode URL encoding (%20, dll.)
		unescapedRel, err := url.PathUnescape(rel)
		if err != nil {
			return fmt.Errorf("failed to decode uri: %w", err)
		}

		localPath := filepath.Join(s.BasePath, filepath.FromSlash(unescapedRel))
		return os.Remove(localPath)
	}

	unescapedURI, err := url.PathUnescape(uri)
	if err != nil {
		return fmt.Errorf("failed to decode uri: %w", err)
	}

	localPath := filepath.Join(s.BasePath, unescapedURI)
	return os.Remove(localPath)
}
