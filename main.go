// VoiceLine — turns sales call audio into structured CRM data.
// Flow: audio upload → Gemini analyzes → row appended to Google Sheet.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"voiceline/internal/config"
	"voiceline/internal/handlers"
	"voiceline/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	geminiSvc, err := services.NewGeminiService(cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("gemini: %v", err)
	}

	sheetsSvc, err := services.NewSheetsService(context.Background(), cfg.SpreadsheetID)
	if err != nil {
		log.Fatalf("sheets: %v", err)
	}

	voiceHandler := &handlers.VoiceHandler{
		Gemini: geminiSvc,
		Sheets: sheetsSvc,
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.StaticFile("/", "web/index.html")
	router.POST("/voice-to-crm", voiceHandler.VoiceToCRM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("VoiceLine server listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
