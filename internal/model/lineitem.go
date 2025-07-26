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
	Name         string   `json:"name" validate:"required,min=1,max=100"`
	AdvertiserID string   `json:"advertiser_id" validate:"required"`
	Bid          float64  `json:"bid" validate:"required,gte=0.1,lte=10"`
	Budget       float64  `json:"budget" validate:"required,gte=1000,lte=10000"`
	Placement    string   `json:"placement" validate:"required,oneof=homepage_sidebar video_preroll article_inline_1 mobile_sticky footer_banner homepage_top article_inline_2"`
	Categories   []string `json:"categories,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`
}
