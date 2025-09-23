package models

import (
	"time"
)

// Bid represents a bid from the auction system
type Bid struct {
	ID           string    `json:"id"`
	EventID      string    `json:"event_id"`
	Round        int       `json:"round"`
	EaselNumber  int       `json:"easel_number"`
	BidderID     string    `json:"bidder_id"`
	Amount       float64   `json:"amount"`
	IsWinning    bool      `json:"is_winning"`
	BidTime      time.Time `json:"bid_time"`
	PaymentStatus string   `json:"payment_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Joined bidder information
	BidderName  string `json:"bidder_name,omitempty"`
	BidderEmail string `json:"bidder_email,omitempty"`
	BidderPhone string `json:"bidder_phone,omitempty"`
}

// AuctionLot represents aggregated auction data for an easel
type AuctionLot struct {
	EventID      string  `json:"event_id"`
	Round        int     `json:"round"`
	EaselNumber  int     `json:"easel_number"`
	ArtistName   string  `json:"artist_name"`
	BidCount     int     `json:"bid_count"`
	HighestBid   float64 `json:"highest_bid"`
	WinningBid   *Bid    `json:"winning_bid,omitempty"`
	AllBids      []Bid   `json:"all_bids,omitempty"`
}