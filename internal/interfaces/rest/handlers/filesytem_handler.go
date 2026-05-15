package handlers

import (
	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/biacibengamukulu/stream-file-locally/internal/interfaces/rest/controllers"
	"github.com/gofiber/fiber/v2"
)

type FileSystemHandler struct {
	svc        filesystem.FileService
	controller controllers.FilesystemController
}

func NewFileSystemHandler(svc filesystem.FileService) *FileSystemHandler {
	controller := controllers.NewFilesystemController(svc)
	return &FileSystemHandler{svc: svc, controller: controller}
}
func (h *FileSystemHandler) Register(router fiber.Router) {
	controller := h.controller
	router.Post("/upload", controller.Upload)
	router.Get("/:id", controller.Get)
}
