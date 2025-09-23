package database

import (
	"context"
	"fmt"
	"paperwork-service/internal/config"
	"paperwork-service/internal/models"

	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// SupabaseClient wraps the Supabase client with our business logic
type SupabaseClient struct {
	client *supabase.Client
	logger *zap.Logger
	config *config.Config
}

// NewSupabaseClient creates a new Supabase client wrapper
func NewSupabaseClient(cfg *config.Config, logger *zap.Logger) (*SupabaseClient, error) {
	client, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &SupabaseClient{
		client: client,
		logger: logger,
		config: cfg,
	}, nil
}

// GetEventByEID fetches an event by its EID (e.g., "AB2995")
func (s *SupabaseClient) GetEventByEID(ctx context.Context, eid string) (*models.Event, error) {
	s.logger.Info("Fetching event by EID", zap.String("eid", eid))

	var events []models.Event
	err := s.client.DB.From("events").
		Select("*").
		Eq("eid", eid).
		Execute(&events)

	if err != nil {
		s.logger.Error("Failed to fetch event", zap.String("eid", eid), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch event %s: %w", eid, err)
	}

	if len(events) == 0 {
		s.logger.Warn("Event not found", zap.String("eid", eid))
		return nil, fmt.Errorf("event not found: %s", eid)
	}

	return &events[0], nil
}

// GetEventArtists fetches all artists for an event with their profile information
func (s *SupabaseClient) GetEventArtists(ctx context.Context, eventID string) ([]models.EventArtist, error) {
	s.logger.Info("Fetching event artists", zap.String("event_id", eventID))

	query := `
		*,
		person:people!inner(
			id,
			email,
			phone,
			name,
			first_name,
			last_name,
			nickname,
			display_phone
		),
		artist_profile:artist_profiles(
			id,
			artist_name,
			bio,
			website,
			instagram,
			facebook,
			twitter,
			tiktok,
			youtube,
			profile_image_url
		)
	`

	var contestants []struct {
		models.EventArtist
		Person        models.Person        `json:"person"`
		ArtistProfile models.ArtistProfile `json:"artist_profile"`
	}

	err := s.client.DB.From("round_contestants").
		Select(query).
		Eq("event_id", eventID).
		Order("round", &supabase.OrderOpts{Ascending: true}).
		Order("easel_number", &supabase.OrderOpts{Ascending: true}).
		Execute(&contestants)

	if err != nil {
		s.logger.Error("Failed to fetch event artists", zap.String("event_id", eventID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch event artists: %w", err)
	}

	// Transform the data into our EventArtist model
	artists := make([]models.EventArtist, len(contestants))
	for i, contestant := range contestants {
		artist := contestant.EventArtist
		artist.Person = contestant.Person
		artist.ArtistProfile = contestant.ArtistProfile

		// Map fields for easier access in templates
		artist.FirstName = contestant.Person.FirstName
		artist.LastName = contestant.Person.LastName
		artist.Email = contestant.Person.Email
		artist.Phone = contestant.Person.Phone
		artist.Instagram = contestant.ArtistProfile.Instagram
		artist.Bio = contestant.ArtistProfile.Bio

		// TODO: Add city/province/country mapping from related tables if needed
		artists[i] = artist
	}

	s.logger.Info("Successfully fetched event artists",
		zap.String("event_id", eventID),
		zap.Int("count", len(artists)))

	return artists, nil
}

// GetEventBids fetches all bids for an event with bidder information
func (s *SupabaseClient) GetEventBids(ctx context.Context, eventID string) ([]models.Bid, error) {
	s.logger.Info("Fetching event bids", zap.String("event_id", eventID))

	query := `
		*,
		bidder:people!inner(
			id,
			name,
			email,
			phone,
			first_name,
			last_name
		)
	`

	var bidsData []struct {
		models.Bid
		Bidder models.Person `json:"bidder"`
	}

	err := s.client.DB.From("bids").
		Select(query).
		Eq("event_id", eventID).
		Order("round", &supabase.OrderOpts{Ascending: true}).
		Order("easel_number", &supabase.OrderOpts{Ascending: true}).
		Order("amount", &supabase.OrderOpts{Ascending: false}).
		Execute(&bidsData)

	if err != nil {
		s.logger.Error("Failed to fetch event bids", zap.String("event_id", eventID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch event bids: %w", err)
	}

	// Transform the data
	bids := make([]models.Bid, len(bidsData))
	for i, bidData := range bidsData {
		bid := bidData.Bid
		bid.BidderName = bidData.Bidder.Name
		if bid.BidderName == "" {
			bid.BidderName = fmt.Sprintf("%s %s", bidData.Bidder.FirstName, bidData.Bidder.LastName)
		}
		bid.BidderEmail = bidData.Bidder.Email
		bid.BidderPhone = bidData.Bidder.Phone

		bids[i] = bid
	}

	s.logger.Info("Successfully fetched event bids",
		zap.String("event_id", eventID),
		zap.Int("count", len(bids)))

	return bids, nil
}

// GetArtistEventHistory fetches an artist's event participation history
func (s *SupabaseClient) GetArtistEventHistory(ctx context.Context, personID string) ([]models.ArtistEvent, error) {
	s.logger.Info("Fetching artist event history", zap.String("person_id", personID))

	query := `
		round,
		easel_number,
		event:events!inner(
			id,
			eid,
			name,
			event_start_datetime
		)
	`

	var historyData []struct {
		Round       int `json:"round"`
		EaselNumber int `json:"easel_number"`
		Event       struct {
			ID                 string `json:"id"`
			EID                string `json:"eid"`
			Name               string `json:"name"`
			EventStartDatetime string `json:"event_start_datetime"`
		} `json:"event"`
	}

	err := s.client.DB.From("round_contestants").
		Select(query).
		Eq("person_id", personID).
		Order("event.event_start_datetime", &supabase.OrderOpts{Ascending: false}).
		Limit(10). // Limit to last 10 events
		Execute(&historyData)

	if err != nil {
		s.logger.Error("Failed to fetch artist event history", zap.String("person_id", personID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch artist event history: %w", err)
	}

	history := make([]models.ArtistEvent, len(historyData))
	for i, item := range historyData {
		// Parse date string - Supabase returns ISO 8601 format
		// You might need to adjust this parsing based on the actual format
		history[i] = models.ArtistEvent{
			EventID:     item.Event.ID,
			EventEID:    item.Event.EID,
			EventName:   item.Event.Name,
			Round:       item.Round,
			EaselNumber: item.EaselNumber,
			// Date parsing would be handled here if needed
		}
	}

	return history, nil
}