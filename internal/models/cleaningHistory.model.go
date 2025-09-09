package models

import (
	"time"

	"github.com/google/uuid"
)

type CleaningHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"                         json:"userId"`
	User      User      `gorm:"foreignKey:UserID"                                json:"user"`
	ReleaseID uuid.UUID `gorm:"type:uuid;not null;index"                         json:"releaseId"`
	Release   Release   `gorm:"foreignKey:ReleaseID"                             json:"release"`
	CleanedAt time.Time `gorm:"not null"                                         json:"cleanedAt"`
	Notes     string    `                                                        json:"notes"`
	CreatedAt time.Time `gorm:"autoCreateTime"                                   json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"                                   json:"updatedAt"`
}

