package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlayHistory represents a user's play history (user-specific)
type PlayHistory struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	ReleaseID uuid.UUID  `gorm:"type:uuid;not null;index" json:"releaseId"`
	StylusID  *uuid.UUID `gorm:"type:uuid" json:"stylusId"`
	PlayedAt  time.Time  `gorm:"not null" json:"playedAt"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Release Release  `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
	Stylus  *Stylus  `gorm:"foreignKey:StylusID" json:"stylus,omitempty"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (p *PlayHistory) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}