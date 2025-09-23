package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"paperwork-service/internal/models"

	"go.uber.org/zap"
)

// EventService handles fetching event data from Supabase edge functions
type EventService struct {
	logger      *zap.Logger
	supabaseURL string
	httpClient  *http.Client
}

// NewEventService creates a new event service
func NewEventService(logger *zap.Logger, supabaseURL string) *EventService {
	return &EventService{
		logger:      logger,
		supabaseURL: supabaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PaperworkData represents the response from the paperwork-data edge function
type PaperworkData struct {
	Event       models.Event             `json:"event"`
	Artists     []models.EventArtist     `json:"artists"`
	AuctionLots []models.AuctionLot      `json:"auction_lots"`
	TotalArtists int                     `json:"total_artists"`
	TotalBids    int                     `json:"total_bids"`
	GeneratedAt  string                  `json:"generated_at"`
}

// GetEventPaperworkData fetches all data needed for paperwork generation
func (s *EventService) GetEventPaperworkData(ctx context.Context, eid string) (*PaperworkData, error) {
	s.logger.Info("Fetching paperwork data for event", zap.String("eid", eid))

	// Build the edge function URL
	url := fmt.Sprintf("%s/functions/v1/paperwork-data/%s", s.supabaseURL, eid)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to make request to edge function",
			zap.String("url", url),
			zap.Error(err))
		return nil, fmt.Errorf("failed to fetch data from edge function: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Edge function returned error",
			zap.String("eid", eid),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)))

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("event not found: %s", eid)
		}
		return nil, fmt.Errorf("edge function error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var data PaperworkData
	if err := json.Unmarshal(body, &data); err != nil {
		s.logger.Error("Failed to parse response JSON",
			zap.String("eid", eid),
			zap.Error(err),
			zap.String("response", string(body)))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	s.logger.Info("Successfully fetched paperwork data",
		zap.String("eid", eid),
		zap.String("event_name", data.Event.Name),
		zap.Int("total_artists", data.TotalArtists),
		zap.Int("total_bids", data.TotalBids))

	return &data, nil
}

// GetEventByEID is a convenience method to get just the event data
func (s *EventService) GetEventByEID(ctx context.Context, eid string) (*models.Event, error) {
	data, err := s.GetEventPaperworkData(ctx, eid)
	if err != nil {
		return nil, err
	}
	return &data.Event, nil
}