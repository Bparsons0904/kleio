package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRelease represents the many-to-many relationship between users and releases (collection ownership)
type UserRelease struct {
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"userId"`
	ReleaseID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"releaseId"`
	InstanceID int       `json:"instanceId"` // Discogs instance ID
	FolderID   int       `json:"folderId"`   // User's Discogs folder
	Rating     int       `json:"rating"`     // User's personal rating
	Notes      string    `json:"notes"`      // User's personal notes
	AddedAt    time.Time `gorm:"autoCreateTime" json:"addedAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Release Release `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
}