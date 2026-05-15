package filesystem

import "time"

type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	URL         string    `json:"url"`
	Extension   string    `json:"extension"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

type StoredFile struct {
	File
	Content []byte `json:"-"`
}
