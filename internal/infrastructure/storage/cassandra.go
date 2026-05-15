package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/biacibengamukulu/stream-file-locally/sharded/cassandra"
	"github.com/gocql/gocql"
)

const filesTableDDL = `
CREATE TABLE IF NOT EXISTS files (
	id text PRIMARY KEY,
	name text,
	content_type text,
	url text,
	extension text,
	size bigint,
	content blob,
	created_at timestamp
)`

type CassandraStorage struct {
	session *cassandra.Session
}

func NewCassandraStorage(session *cassandra.Session) (*CassandraStorage, error) {
	if session == nil || session.Session == nil {
		return nil, fmt.Errorf("cassandra session is required")
	}
	return &CassandraStorage{session: session}, nil
}

func CassandraDDLs() []string {
	return []string{filesTableDDL}
}

func (s *CassandraStorage) Save(file filesystem.StoredFile) (filesystem.File, error) {
	err := s.session.Query(
		`INSERT INTO files (id, name, content_type, url, extension, size, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		file.ID,
		file.Name,
		file.ContentType,
		file.URL,
		file.Extension,
		file.Size,
		file.Content,
		file.CreatedAt,
	).Exec()
	if err != nil {
		return filesystem.File{}, err
	}
	return file.File, nil
}

func (s *CassandraStorage) Get(id string) (filesystem.File, error) {
	var file filesystem.File
	if err := s.session.Query(
		`SELECT id, name, content_type, url, extension, size, created_at FROM files WHERE id = ?`,
		id,
	).Scan(
		&file.ID,
		&file.Name,
		&file.ContentType,
		&file.URL,
		&file.Extension,
		&file.Size,
		&file.CreatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return filesystem.File{}, filesystem.ErrNotFound
		}
		return filesystem.File{}, err
	}
	return file, nil
}

func (s *CassandraStorage) Read(id string) (filesystem.StoredFile, error) {
	var file filesystem.StoredFile
	var createdAt time.Time
	if err := s.session.Query(
		`SELECT id, name, content_type, url, extension, size, content, created_at FROM files WHERE id = ?`,
		id,
	).Scan(
		&file.ID,
		&file.Name,
		&file.ContentType,
		&file.URL,
		&file.Extension,
		&file.Size,
		&file.Content,
		&createdAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return filesystem.StoredFile{}, filesystem.ErrNotFound
		}
		return filesystem.StoredFile{}, err
	}
	file.CreatedAt = createdAt
	return file, nil
}

func (s *CassandraStorage) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	return s.session.Query(`DELETE FROM files WHERE id = ?`, id).Exec()
}
