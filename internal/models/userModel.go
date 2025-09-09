package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the multi-user system
type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	Username        string    `gorm:"uniqueIndex;not null" json:"username"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	DiscogsUsername string    `json:"discogsUsername"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	AuthTokens      []AuthToken       `gorm:"foreignKey:UserID" json:"-"`
	UserReleases    []UserRelease     `gorm:"foreignKey:UserID" json:"-"`
	PlayHistory     []PlayHistory     `gorm:"foreignKey:UserID" json:"-"`
	Styluses        []Stylus          `gorm:"foreignKey:UserID" json:"-"`
	CleaningHistory []CleaningHistory `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook for generating UUIDs if not provided
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}