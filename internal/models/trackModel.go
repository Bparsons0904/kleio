package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Track represents a track on a release (shared across users)
type Track struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	ReleaseID       uuid.UUID `gorm:"type:uuid;not null;index" json:"releaseId"`
	Position        string    `gorm:"not null" json:"position"` // e.g., "A1", "B2"
	Title           string    `gorm:"not null" json:"title"`
	DurationText    string    `json:"durationText"`    // Original format (e.g., "3:45")
	DurationSeconds int       `json:"durationSeconds"` // Normalized duration
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Release Release `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (t *Track) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}