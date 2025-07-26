package model

import "time"

// TrackingEventType represents the type of tracking event
type TrackingEventType string

const (
	TrackingEventTypeImpression TrackingEventType = "impression"
	TrackingEventTypeClick      TrackingEventType = "click"
	TrackingEventTypeConversion TrackingEventType = "conversion"
)

// TrackingEvent represents a user interaction with an ad
type TrackingEvent struct {
	EventType  TrackingEventType `json:"event_type" query:"event_type" validate:"required,oneof=click conversion impression"`
	LineItemID string            `json:"line_item_id" query:"line_item_id" validate:"required"`
	Timestamp  time.Time         `json:"timestamp,omitempty" query:"timestamp" validate:"omitempty"`
	Placement  string            `json:"placement,omitempty" query:"placement" validate:"required"`
	UserID     string            `json:"user_id,omitempty" query:"user_id" validate:"required"`
	Metadata   map[string]string `json:"metadata,omitempty" query:"metadata" validate:"required"`
}
