package models

import (
	"time"
)

// Person represents a person/user from Supabase
type Person struct {
	ID                   string    `json:"id"`
	Email                string    `json:"email"`
	Phone                string    `json:"phone"`
	Name                 string    `json:"name"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Nickname             string    `json:"nickname"`
	DisplayPhone         string    `json:"display_phone"`
	RegionCode           string    `json:"region_code"`
	LastInteractionAt    time.Time `json:"last_interaction_at"`
	InteractionCount     int       `json:"interaction_count"`
	TotalSpent           float64   `json:"total_spent"`
	LastQrScanAt         time.Time `json:"last_qr_scan_at"`
	LastQrEventID        string    `json:"last_qr_event_id"`
	ArtBattleNews        bool      `json:"art_battle_news"`
	NotificationEmails   bool      `json:"notification_emails"`
	LoyaltyOffers        bool      `json:"loyalty_offers"`
	VerificationCode     string    `json:"verification_code"`
	VerificationCodeExp  time.Time `json:"verification_code_exp"`
	SelfRegistered       bool      `json:"self_registered"`
	MessageBlocked       int       `json:"message_blocked"`
	LocationLat          float64   `json:"location_lat"`
	LocationLng          float64   `json:"location_lng"`
	RegisteredAt         string    `json:"registered_at"`
	LastPromoSentAt      time.Time `json:"last_promo_sent_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ArtistProfile represents an artist profile from Supabase
type ArtistProfile struct {
	ID                   string    `json:"id"`
	PersonID             string    `json:"person_id"`
	ArtistName           string    `json:"artist_name"`
	Bio                  string    `json:"bio"`
	Website              string    `json:"website"`
	Instagram            string    `json:"instagram"`
	Facebook             string    `json:"facebook"`
	Twitter              string    `json:"twitter"`
	TikTok               string    `json:"tiktok"`
	YouTube              string    `json:"youtube"`
	ProfileImageURL      string    `json:"profile_image_url"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// EventArtist represents an artist participating in an event (matches edge function response)
type EventArtist struct {
	ContestantID     string `json:"contestant_id"`
	EaselNumber      int    `json:"easel_number"`
	RoundNumber      int    `json:"round_number"`
	EventID          string `json:"event_id"`
	ArtistProfileID  string `json:"artist_profile_id"`
	EntryID          int    `json:"entry_id"`
	ArtistName       string `json:"artist_name"`
	Bio              string `json:"bio"`
	ABHQBio          string `json:"abhq_bio"`
	Instagram        string `json:"instagram"`
	Website          string `json:"website"`
	PersonName       string `json:"person_name"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	DisplayName      string `json:"display_name"`
	Status           string `json:"status"`

	// Legacy fields for compatibility with PDF generation
	Round           int           `json:"round"`
	EventHistory    []ArtistEvent `json:"event_history,omitempty"`
}

// ArtistEvent represents an artist's participation in past events
type ArtistEvent struct {
	EventID     string    `json:"event_id"`
	EventName   string    `json:"event_name"`
	EventEID    string    `json:"event_eid"`
	EventDate   time.Time `json:"event_date"`
	Round       int       `json:"round_number"`
	EaselNumber int       `json:"easel_number"`
	IsWinner    bool      `json:"is_winner"`
}