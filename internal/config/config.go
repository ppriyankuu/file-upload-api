package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          int
	UploadDir     string
	MaxUploadSize int64 // in bytes
	AllowedTypes  []string
}

// Load reads from envs
func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	max := int64(10 << 20) // 10 mb default
	if m := os.Getenv("MAX_UPLOAD_SIZE"); m != "" {
		if v, err := strconv.ParseInt(m, 10, 64); err == nil {
			max = v
		}
	}

	allowed := []string{"image/jpeg", "image/png", "application/pdf"}
	if a := os.Getenv("ALLOWED_TYPES"); a != "" {
		allowed = parseCSV(a)
	}

	cfg := &Config{
		Port:          port,
		UploadDir:     uploadDir,
		MaxUploadSize: max,
		AllowedTypes:  allowed,
	}

	log.Printf("config: port=%d upload_dir=%s max_upload_size=%d allowed=%v", cfg.Port, cfg.UploadDir, cfg.MaxUploadSize, cfg.AllowedTypes)
	return cfg
}

func parseCSV(s string) []string {
	out := []string{}
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
