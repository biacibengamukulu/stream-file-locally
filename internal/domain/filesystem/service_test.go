package filesystem_test

import (
	"testing"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/biacibengamukulu/stream-file-locally/internal/infrastructure/storage"
	"github.com/biacibengamukulu/stream-file-locally/sharded/config"
)

func TestFileServiceUploadAndStreamWithDiskStorage(t *testing.T) {
	disk, err := storage.NewDiskStorage(t.TempDir())
	if err != nil {
		t.Fatalf("new disk storage: %v", err)
	}

	cfg := config.Config{
		RoutePrefix:   "/stream-file-locally",
		PublicBaseURL: "http://example.test",
	}
	svc := filesystem.NewFileService(cfg, disk)

	uploaded, err := svc.Upload(filesystem.FileDTO{
		Name:          "hello.txt",
		ContentType:   "text/plain",
		Base64Content: "SGVsbG8sIHdvcmxkIQ==",
	})
	if err != nil {
		t.Fatalf("upload: %v", err)
	}

	if uploaded.ID == "" {
		t.Fatal("expected generated id")
	}
	if uploaded.URL == "" {
		t.Fatal("expected stream url")
	}
	if uploaded.Size != 13 {
		t.Fatalf("expected size 13, got %d", uploaded.Size)
	}

	meta, err := svc.Get(uploaded.ID)
	if err != nil {
		t.Fatalf("get metadata: %v", err)
	}
	if meta.Name != "hello.txt" {
		t.Fatalf("expected hello.txt, got %q", meta.Name)
	}

	streamed, err := svc.Stream(uploaded.ID)
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if string(streamed.Content) != "Hello, world!" {
		t.Fatalf("unexpected stream content: %q", string(streamed.Content))
	}

	if err := svc.Delete(uploaded.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(uploaded.ID); err != filesystem.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestFileServiceUploadAcceptsDataURLBase64(t *testing.T) {
	disk, err := storage.NewDiskStorage(t.TempDir())
	if err != nil {
		t.Fatalf("new disk storage: %v", err)
	}

	svc := filesystem.NewFileService(config.Config{
		RoutePrefix: "/stream-file-locally",
	}, disk)

	uploaded, err := svc.Upload(filesystem.FileDTO{
		Name:          "hello.txt",
		Base64Content: "data:text/plain;base64,SGVsbG8sIGRhdGEgVVJMIQ==",
	})
	if err != nil {
		t.Fatalf("upload data url: %v", err)
	}
	if uploaded.ContentType != "text/plain" {
		t.Fatalf("expected content type from data url, got %q", uploaded.ContentType)
	}

	streamed, err := svc.Stream(uploaded.ID)
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if string(streamed.Content) != "Hello, data URL!" {
		t.Fatalf("unexpected stream content: %q", string(streamed.Content))
	}
}
