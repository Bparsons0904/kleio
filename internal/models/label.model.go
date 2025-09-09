package models

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
	DiscogsID   int       `gorm:"uniqueIndex;not null"                             json:"discogsId"`
	Name        string    `gorm:"not null"                                         json:"name"`
	ResourceURL string    `                                                        json:"resourceUrl"`
	EntityType  string    `                                                        json:"entityType"`
	CreatedAt   time.Time `gorm:"autoCreateTime"                                   json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"                                   json:"updatedAt"`
	Releases    []Release `gorm:"many2many:release_labels"                         json:"releases"`
}

