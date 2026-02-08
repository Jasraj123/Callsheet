package models

// VoiceAnalysis is the structured output from Gemini for a sales audio.
type VoiceAnalysis struct {
	Summary      string   `json:"summary"`
	ActionItems  []string `json:"action_items"`
	Sentiment    string   `json:"sentiment"` // positive, neutral, negative
	UrgencyScore int      `json:"urgency_score"`
	ClientName   string   `json:"client_name"`
}
