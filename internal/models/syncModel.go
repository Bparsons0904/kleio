package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sync represents synchronization status (user-specific)
type Sync struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	SyncStart time.Time  `json:"syncStart"`
	SyncEnd   *time.Time `json:"syncEnd,omitempty"`
	Status    string     `json:"status"` // "in_progress" or "complete" or "failed"
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (s *Sync) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}