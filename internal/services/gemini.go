package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"google.golang.org/genai"

	"voiceline/internal/models"
)

const (
	geminiModel     = "gemini-2.5-flash"
	analysisTimeout = 60 * time.Second
)

var (
	salesPrompt = `Act as a Sales Assistant. Analyze this audio and return a JSON object with exactly these fields (no markdown, no code block, only valid JSON):
- summary (string): brief summary of the conversation
- action_items (array of strings): list of follow-up actions
- sentiment (string): one of "positive", "neutral", "negative"
- urgency_score (number): integer from 1 to 10
- client_name (string): name of the client or contact mentioned

Return only the JSON object, nothing else.`
)

type GeminiService struct {
	client *genai.Client
}

// NewGeminiService creates a Gemini service with the given API key.
func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("genai client: %w", err)
	}
	return &GeminiService{client: client}, nil
}

// AnalyzeAudio uploads the audio file via Files API, sends it to Gemini with the sales prompt, and returns structured analysis.
func (s *GeminiService) AnalyzeAudio(ctx context.Context, audioPath, mimeType string) (*models.VoiceAnalysis, error) {
	ctx, cancel := context.WithTimeout(ctx, analysisTimeout)
	defer cancel()

	uploadedFile, err := s.client.Files.UploadFromPath(ctx, audioPath, &genai.UploadFileConfig{MIMEType: mimeType})
	if err != nil {
		return nil, fmt.Errorf("upload audio: %w", err)
	}
	defer func() { _, _ = s.client.Files.Delete(ctx, uploadedFile.Name, nil) }()

	parts := []*genai.Part{
		genai.NewPartFromText(salesPrompt),
		genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
	}

	contents := []*genai.Content{genai.NewContentFromParts(parts, genai.RoleUser)}
	resp, err := s.client.Models.GenerateContent(ctx, geminiModel, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("generate content: %w", err)
	}

	text := resp.Text()
	if text == "" {
		return nil, fmt.Errorf("empty response from model")
	}

	analysis, err := parseAnalysisJSON(text)
	if err != nil {
		return nil, fmt.Errorf("parse analysis: %w", err)
	}

	return analysis, nil
}

// parseAnalysisJSON extracts and unmarshals JSON from model output (may be wrapped in markdown).
func parseAnalysisJSON(text string) (*models.VoiceAnalysis, error) {
	text = trimMarkdownJSON(text)
	var out models.VoiceAnalysis
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return nil, err
	}
	// Clamp urgency to 1-10
	if out.UrgencyScore < 1 {
		out.UrgencyScore = 1
	}
	if out.UrgencyScore > 10 {
		out.UrgencyScore = 10
	}
	return &out, nil
}

var jsonBlockRe = regexp.MustCompile("(?s)`{0,3}json\\s*\\n?(.*?)`{0,3}\\s*$")

func trimMarkdownJSON(s string) string {
	if match := jsonBlockRe.FindStringSubmatch(s); len(match) > 1 {
		return match[1]
	}
	return s
}
