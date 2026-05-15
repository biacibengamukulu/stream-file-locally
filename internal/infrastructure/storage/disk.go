package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
)

type DiskStorage struct {
	root string
}

func NewDiskStorage(root string) (*DiskStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("disk storage path is required")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &DiskStorage{root: root}, nil
}

func (s *DiskStorage) Save(file filesystem.StoredFile) (filesystem.File, error) {
	if file.ID == "" {
		return filesystem.File{}, fmt.Errorf("file id is required")
	}
	dir := filepath.Join(s.root, file.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return filesystem.File{}, err
	}

	contentPath := filepath.Join(dir, "content")
	if err := os.WriteFile(contentPath, file.Content, 0o644); err != nil {
		return filesystem.File{}, err
	}

	metaPath := filepath.Join(dir, "metadata.json")
	meta, err := json.MarshalIndent(file.File, "", "  ")
	if err != nil {
		return filesystem.File{}, err
	}
	if err := os.WriteFile(metaPath, meta, 0o644); err != nil {
		return filesystem.File{}, err
	}

	return file.File, nil
}

func (s *DiskStorage) Get(id string) (filesystem.File, error) {
	var file filesystem.File
	metaPath := filepath.Join(s.root, id, "metadata.json")
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return filesystem.File{}, filesystem.ErrNotFound
		}
		return filesystem.File{}, err
	}
	if err := json.Unmarshal(raw, &file); err != nil {
		return filesystem.File{}, err
	}
	return file, nil
}

func (s *DiskStorage) Read(id string) (filesystem.StoredFile, error) {
	file, err := s.Get(id)
	if err != nil {
		return filesystem.StoredFile{}, err
	}

	contentPath := filepath.Join(s.root, id, "content")
	content, err := os.ReadFile(contentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return filesystem.StoredFile{}, filesystem.ErrNotFound
		}
		return filesystem.StoredFile{}, err
	}

	return filesystem.StoredFile{File: file, Content: content}, nil
}

func (s *DiskStorage) Delete(id string) error {
	dir := filepath.Join(s.root, id)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return filesystem.ErrNotFound
		}
		return err
	}
	return os.RemoveAll(dir)
}
