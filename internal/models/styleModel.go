package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Style represents a music style (shared across users)
type Style struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	Name string    `gorm:"uniqueIndex;not null" json:"name"`

	Releases []Release `gorm:"many2many:release_styles" json:"-"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (s *Style) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}