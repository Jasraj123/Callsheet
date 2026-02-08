// services/sheets — appends analysis rows to a Google Sheet using the Sheets API.
package services

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"voiceline/internal/models"
)

const (
	sheetsScope    = "https://www.googleapis.com/auth/spreadsheets"
	sheetsTimeout  = 15 * time.Second
	defaultSheet   = "Sheet1"
	appendRange    = "Sheet1!A:F"
	valueInputUser = "USER_ENTERED"
)

// SheetsService appends rows to a Google Sheet.
type SheetsService struct {
	svc           *sheets.Service
	spreadsheetID string
}

// NewSheetsService creates a Sheets service using application default credentials
// (e.g. GOOGLE_APPLICATION_CREDENTIALS pointing to a service account JSON).
func NewSheetsService(ctx context.Context, spreadsheetID string) (*SheetsService, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	creds, err := google.FindDefaultCredentials(ctx, sheetsScope)
	if err != nil {
		return nil, fmt.Errorf("find default credentials: %w", err)
	}

	client, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("sheets client: %w", err)
	}

	return &SheetsService{svc: client, spreadsheetID: spreadsheetID}, nil
}

// AppendAnalysisRow adds one row with: Timestamp | Client | Summary | Sentiment | Urgency | Urgent (Yes if urgency > 7).
func (s *SheetsService) AppendAnalysisRow(ctx context.Context, analysis *models.VoiceAnalysis) error {
	ctx, cancel := context.WithTimeout(ctx, sheetsTimeout)
	defer cancel()

	ts := time.Now().Format(time.RFC3339)
	urgent := "No"
	if analysis.UrgencyScore > 7 {
		urgent = "Yes"
	}
	row := []interface{}{
		ts,
		analysis.ClientName,
		analysis.Summary,
		analysis.Sentiment,
		analysis.UrgencyScore,
		urgent,
	}

	_, err := s.svc.Spreadsheets.Values.Append(
		s.spreadsheetID,
		appendRange,
		&sheets.ValueRange{Values: [][]interface{}{row}},
	).ValueInputOption(valueInputUser).InsertDataOption("INSERT_ROWS").Context(ctx).Do()

	if err != nil {
		return fmt.Errorf("sheets append: %w", err)
	}
	return nil
}
