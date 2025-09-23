package models

import (
	"time"
)

// Event represents an Art Battle event from Supabase
type Event struct {
	ID                   string    `json:"id"`
	EID                  string    `json:"eid"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Venue                string    `json:"venue"`
	EventStartDatetime   time.Time `json:"event_start_datetime"`
	EventEndDatetime     time.Time `json:"event_end_datetime"`
	TimezoneID           string    `json:"timezone_id"`
	TimezoneOffset       string    `json:"timezone_offset"`
	TimezoneIcann        string    `json:"timezone_icann"`
	CityID               string    `json:"city_id"`
	CountryID            string    `json:"country_id"`
	Enabled              bool      `json:"enabled"`
	ShowInApp            bool      `json:"show_in_app"`
	CurrentRound         int       `json:"current_round"`
	ArtWidthHeight       string    `json:"art_width_height"`
	VoteByLink           bool      `json:"vote_by_link"`
	RegisterAtSmsVote    bool      `json:"register_at_sms_vote"`
	SendLinkToGuests     bool      `json:"send_link_to_guests"`
	EmailRegistration    bool      `json:"email_registration"`
	EnableAuction        bool      `json:"enable_auction"`
	AuctionStartBid      float64   `json:"auction_start_bid"`
	MinBidIncrement      float64   `json:"min_bid_increment"`
	Currency             string    `json:"currency"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// EventWithDetails includes related data for paperwork generation
type EventWithDetails struct {
	Event
	Artists []EventArtist `json:"artists"`
	Bids    []Bid         `json:"bids"`
}