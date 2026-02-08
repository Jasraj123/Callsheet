package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration from environment.
type Config struct {
	GeminiAPIKey      string
	SpreadsheetID     string
	GoogleCredentials string 
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		GeminiAPIKey:      strings.TrimSpace(os.Getenv("GEMINI_API_KEY")),
		SpreadsheetID:     strings.TrimSpace(os.Getenv("SPREADSHEET_ID")),
		GoogleCredentials: strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")),
	}

	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}
	if cfg.SpreadsheetID == "" {
		return nil, fmt.Errorf("SPREADSHEET_ID is required")
	}

	return cfg, nil
}
