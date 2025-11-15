package main

import (
	"file-upload-api/internal/config"
	"file-upload-api/internal/handlers"
	"file-upload-api/internal/middleware"
	"file-upload-api/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// ensuring upload dir exists
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}

	store := storage.NewLocalStorage(cfg.UploadDir)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.LimitRequestBody(cfg.MaxUploadSize))

	r.StaticFS("/uploads", http.Dir(cfg.UploadDir))

	h := handlers.NewUploadHandler(store, cfg)

	r.POST("/upload", h.Upload)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
