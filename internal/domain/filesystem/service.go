package filesystem

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/biacibengamukulu/stream-file-locally/sharded/config"
	"github.com/google/uuid"
)

type FileService interface {
	Upload(file FileDTO) (File, error)
	Get(id string) (File, error)
	Stream(id string) (StoredFile, error)
	Delete(id string) error
}

type FileServiceImpl struct {
	cfg     config.Config
	storage Storage
}

func NewFileService(cfg config.Config, storage Storage) *FileServiceImpl {
	return &FileServiceImpl{cfg: cfg, storage: storage}
}

func (f FileServiceImpl) Upload(file FileDTO) (File, error) {
	if f.storage == nil {
		return File{}, fmt.Errorf("file storage is not configured")
	}

	name := strings.TrimSpace(file.Name)
	if name == "" {
		name = "upload"
	}

	base64Content, dataURLContentType := normalizeBase64Content(file.Base64Content)
	content, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return File{}, fmt.Errorf("invalid base64 content")
	}
	if len(content) == 0 {
		return File{}, fmt.Errorf("file content is empty")
	}

	contentType := strings.TrimSpace(file.ContentType)
	if contentType == "" {
		contentType = dataURLContentType
	}
	if contentType == "" {
		contentType = http.DetectContentType(content)
	}

	extension := strings.ToLower(filepath.Ext(name))
	if extension == "" {
		if exts, err := mime.ExtensionsByType(contentType); err == nil && len(exts) > 0 {
			extension = exts[0]
		}
	}

	id := uuid.NewString()
	stored := StoredFile{
		File: File{
			ID:          id,
			Name:        name,
			ContentType: contentType,
			URL:         f.streamURL(id),
			Extension:   extension,
			Size:        int64(len(content)),
			CreatedAt:   time.Now().UTC(),
		},
		Content: content,
	}

	return f.storage.Save(stored)
}

func (f FileServiceImpl) Get(id string) (File, error) {
	if f.storage == nil {
		return File{}, fmt.Errorf("file storage is not configured")
	}
	file, err := f.storage.Get(strings.TrimSpace(id))
	if err != nil {
		return File{}, err
	}
	file.URL = f.streamURL(file.ID)
	return file, nil
}

func (f FileServiceImpl) Stream(id string) (StoredFile, error) {
	if f.storage == nil {
		return StoredFile{}, fmt.Errorf("file storage is not configured")
	}
	file, err := f.storage.Read(strings.TrimSpace(id))
	if err != nil {
		return StoredFile{}, err
	}
	file.URL = f.streamURL(file.ID)
	return file, nil
}

func (f FileServiceImpl) Delete(id string) error {
	if f.storage == nil {
		return fmt.Errorf("file storage is not configured")
	}
	return f.storage.Delete(strings.TrimSpace(id))
}

func (f FileServiceImpl) streamURL(id string) string {
	path := strings.TrimRight(f.cfg.RoutePrefix, "/") + "/api/v1/files/" + id + "/stream"
	baseURL := strings.TrimRight(f.cfg.PublicBaseURL, "/")
	if baseURL == "" {
		return path
	}
	return baseURL + path
}

func normalizeBase64Content(value string) (string, string) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(strings.ToLower(value), "data:") {
		return value, ""
	}

	header, payload, ok := strings.Cut(value, ",")
	if !ok {
		return value, ""
	}

	contentType := strings.TrimPrefix(header, "data:")
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = contentType[:idx]
	}
	return strings.TrimSpace(payload), strings.TrimSpace(contentType)
}
