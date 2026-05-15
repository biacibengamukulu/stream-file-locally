package controllers

import (
	"net/http"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/gofiber/fiber/v2"
)

type FilesystemController interface {
	Upload(h *fiber.Ctx) error
	Get(h *fiber.Ctx) error
}

type FilesystemControllerImpl struct {
	svc filesystem.FileService
}

func NewFilesystemController(svc filesystem.FileService) *FilesystemControllerImpl {
	return &FilesystemControllerImpl{svc: svc}
}

func (c *FilesystemControllerImpl) Upload(h *fiber.Ctx) error {
	var req filesystem.FileDTO
	if err := h.BodyParser(&req); err != nil {
		return h.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}

	if req.Base64Content == "" {
		return h.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "base64_content_required"})
	}
	fileInfo, err := c.svc.Upload(req)
	if err != nil {
		return h.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.Status(http.StatusCreated).JSON(fiber.Map{"data": fileInfo})
}

func (c *FilesystemControllerImpl) Get(h *fiber.Ctx) error {
	fileId := h.Params("id")
	if fileId == "" {
		return h.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_id_required"})
	}

	data, err := c.svc.Get(fileId)
	if err != nil {
		return h.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.JSON(fiber.Map{"data": data})
}
