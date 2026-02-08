// handlers — HTTP layer. Receives audio, calls Gemini + Sheets, returns JSON.
package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"voiceline/internal/models"
	"voiceline/internal/services"
)

const (
	maxAudioSize   = 25 << 20 // 25 MB
	multipartFile  = "file"
	allowedWav     = "audio/wav"
	allowedMp3     = "audio/mpeg"
	allowedWebm    = "audio/webm"
)

// VoiceHandler handles voice-to-CRM API.
type VoiceHandler struct {
	Gemini *services.GeminiService
	Sheets *services.SheetsService
}

// VoiceToCRM handles POST /voice-to-crm.
// 1. Validate file (type, size)
// 2. Save to temp file
// 3. Call Gemini → get analysis
// 4. Append row to Sheet
// 5. Return analysis JSON
func (h *VoiceHandler) VoiceToCRM(c *gin.Context) {
	// Step 1: Get and validate the uploaded file
	file, header, err := c.Request.FormFile(multipartFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid file; use form key 'file' with .wav, .mp3, or .webm"})
		return
	}
	defer file.Close()

	if header.Size > maxAudioSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 25MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	var mime string
	switch ext {
	case ".wav":
		mime = allowedWav
	case ".mp3":
		mime = allowedMp3
	case ".webm":
		mime = allowedWebm
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .wav, .mp3, and .webm are supported"})
		return
	}

	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "voiceline-*"+ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp file"})
		return
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save upload"})
		return
	}
	if err := tmpFile.Sync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to flush temp file"})
		return
	}
	audioPath := tmpFile.Name()

	// Step 2: Send audio to Gemini, get structured analysis
	analysis, err := h.Gemini.AnalyzeAudio(c.Request.Context(), audioPath, mime)
	if err != nil {
		log.Printf("Gemini analysis failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("analysis failed: %v", err)})
		return
	}

	// Step 3: Append analysis row to Google Sheet
	if err := h.Sheets.AppendAnalysisRow(c.Request.Context(), analysis); err != nil {
		log.Printf("Sheets append failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    fmt.Sprintf("failed to append to sheet: %v", err),
			"analysis": analysis,
		})
		return
	}

	c.JSON(http.StatusOK, responseFromAnalysis(analysis))
}

func responseFromAnalysis(a *models.VoiceAnalysis) gin.H {
	return gin.H{
		"summary":       a.Summary,
		"action_items":  a.ActionItems,
		"sentiment":     a.Sentiment,
		"urgency_score": a.UrgencyScore,
		"client_name":   a.ClientName,
	}
}
