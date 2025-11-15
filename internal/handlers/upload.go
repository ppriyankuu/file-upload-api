package handlers

import (
	"file-upload-api/internal/config"
	"file-upload-api/internal/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	store *storage.LocalStorage
	cfg   *config.Config
}

func NewUploadHandler(s *storage.LocalStorage, cfg *config.Config) *UploadHandler {
	return &UploadHandler{store: s, cfg: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	// gin will by default parse multipart form when you call FormFile

	// single file
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if fh.Size > h.cfg.MaxUploadSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}

	// reading small prefix to detect content type
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	// reset file pointer - need to reopen to pass stream to storage
	// easiest approach: reopen (FormFile gives header we can re-open)
	f.Close()
	f, err = fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}

	// content type validation
	ok := false
	for _, t := range h.cfg.AllowedTypes {
		if t == contentType {
			ok = true
			break
		}
	}

	if !ok {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": fmt.Sprintf("content type %s not allowed", contentType)})
		return
	}

	savedName, err := h.store.Save(f, fh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	url := fmt.Sprintf("/uploads/%s", savedName)
	c.JSON(http.StatusCreated, gin.H{"url": url})
}
