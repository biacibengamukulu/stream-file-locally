package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/biacibengamukulu/stream-file-locally/internal/domain/filesystem"
	"github.com/gofiber/fiber/v2"
)

type FilesystemController interface {
	Upload(h *fiber.Ctx) error
	Get(h *fiber.Ctx) error
	Stream(h *fiber.Ctx) error
	Delete(h *fiber.Ctx) error
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
		if errors.Is(err, filesystem.ErrNotFound) {
			return h.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file_not_found"})
		}
		return h.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.JSON(fiber.Map{"data": data})
}

func (c *FilesystemControllerImpl) Stream(h *fiber.Ctx) error {
	fileID := h.Params("id")
	if fileID == "" {
		return h.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_id_required"})
	}

	file, err := c.svc.Stream(fileID)
	if err != nil {
		if errors.Is(err, filesystem.ErrNotFound) {
			return h.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file_not_found"})
		}
		return h.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.Set(fiber.HeaderContentType, file.ContentType)
	h.Set(fiber.HeaderContentLength, strconv.FormatInt(int64(len(file.Content)), 10))
	h.Set(fiber.HeaderContentDisposition, `inline; filename="`+safeHeaderFilename(file.Name)+`"`)
	return h.Send(file.Content)
}

func (c *FilesystemControllerImpl) Delete(h *fiber.Ctx) error {
	fileID := h.Params("id")
	if fileID == "" {
		return h.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_id_required"})
	}

	if err := c.svc.Delete(fileID); err != nil {
		if errors.Is(err, filesystem.ErrNotFound) {
			return h.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file_not_found"})
		}
		return h.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.SendStatus(fiber.StatusNoContent)
}

func safeHeaderFilename(name string) string {
	name = strings.ReplaceAll(name, `"`, "")
	name = strings.ReplaceAll(name, "\r", "")
	name = strings.ReplaceAll(name, "\n", "")
	if strings.TrimSpace(name) == "" {
		return "download"
	}
	return name
}
