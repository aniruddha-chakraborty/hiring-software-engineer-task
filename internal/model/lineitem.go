package model

import (
	"time"
)

// LineItemStatus represents the status of a line item
type LineItemStatus string

const (
	LineItemStatusActive    LineItemStatus = "active"
	LineItemStatusPaused    LineItemStatus = "paused"
	LineItemStatusCompleted LineItemStatus = "completed"
)

// LineItem represents an advertisement with associated bid information
type LineItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	AdvertiserID string         `json:"advertiser_id"`
	Bid          float64        `json:"bid"`
	Budget       float64        `json:"budget"`
	Placement    string         `json:"placement"`
	Categories   []string       `json:"categories,omitempty"`
	Keywords     []string       `json:"keywords,omitempty"`
	Status       LineItemStatus `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// LineItemCreate represents the data needed to create a new line item
type LineItemCreate struct {
	Name         string   `json:"name"`
	AdvertiserID string   `json:"advertiser_id"`
	Bid          float64  `json:"bid"`
	Budget       float64  `json:"budget"`
	Placement    string   `json:"placement"`
	Categories   []string `json:"categories,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`
}

// Ad represents an advertisement ready to be served
type Ad struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	AdvertiserID string  `json:"advertiser_id"`
	Bid          float64 `json:"bid"`
	Placement    string  `json:"placement"`
	ServeURL     string  `json:"serve_url"`
	Relevance    int     `json:"relevance"`
}

// TrackingEventType represents the type of tracking event
type TrackingEventType string

const (
	TrackingEventTypeImpression TrackingEventType = "impression"
	TrackingEventTypeClick      TrackingEventType = "click"
	TrackingEventTypeConversion TrackingEventType = "conversion"
)

// TrackingEvent represents a user interaction with an ad
type TrackingEvent struct {
	EventType  TrackingEventType `json:"event_type"`
	LineItemID string            `json:"line_item_id"`
	Timestamp  time.Time         `json:"timestamp,omitempty"`
	Placement  string            `json:"placement,omitempty"`
	UserID     string            `json:"user_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
