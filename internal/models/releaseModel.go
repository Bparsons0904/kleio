package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Release represents a vinyl release (shared across users)
type Release struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	DiscogsID             int       `gorm:"uniqueIndex;not null" json:"discogsId"`
	Title                 string    `gorm:"not null" json:"title"`
	Year                  *int      `json:"year"`
	ResourceURL           string    `json:"resourceUrl"`
	Thumb                 string    `json:"thumb"`
	CoverImage            string    `json:"coverImage"`
	PlayDuration          *int      `json:"playDuration"`
	PlayDurationEstimated *bool     `json:"playDurationEstimated"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Artists      []Artist      `gorm:"many2many:release_artists" json:"artists,omitempty"`
	Labels       []Label       `gorm:"many2many:release_labels" json:"labels,omitempty"`
	Genres       []Genre       `gorm:"many2many:release_genres" json:"genres,omitempty"`
	Styles       []Style       `gorm:"many2many:release_styles" json:"styles,omitempty"`
	Tracks       []Track       `gorm:"foreignKey:ReleaseID" json:"tracks,omitempty"`
	UserReleases []UserRelease `gorm:"foreignKey:ReleaseID" json:"-"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (r *Release) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}