package models

import (
	"github.com/google/uuid"
)

type Genre struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	Name     string    `gorm:"uniqueIndex;not null"                             json:"name"`
	Releases []Release `gorm:"many2many:release_genres"                         json:"releases"`
}

