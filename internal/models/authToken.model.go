package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	Token        string     `gorm:"primaryKey"               json:"token"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	User         User       `gorm:"foreignKey:UserID"        json:"user"`
	DiscogsToken string     `                                json:"discogsToken"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"           json:"createdAt"`
	ExpiresAt    *time.Time `                                json:"expiresAt"`
}

