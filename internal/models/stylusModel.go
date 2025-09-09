package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Stylus represents a stylus owned by a user (user-specific)
type Stylus struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	Name             string     `gorm:"not null" json:"name"`
	Manufacturer     string     `json:"manufacturer"`
	ExpectedLifespan int        `json:"expectedLifespan"`
	PurchaseDate     *time.Time `json:"purchaseDate"`
	Active           bool       `gorm:"default:false" json:"active"`
	Primary          bool       `gorm:"default:false" json:"primary"`
	Owned            bool       `gorm:"default:false" json:"owned"`
	BaseModel        bool       `gorm:"default:false" json:"baseModel"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (s *Stylus) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}