package models

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	ID          int       `gorm:"primaryKey"               json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User        User      `gorm:"foreignKey:UserID"        json:"user"`
	Name        string    `gorm:"not null"                 json:"name"`
	Count       int       `                                json:"count"`
	ResourceURL string    `gorm:"varchar(255)"             json:"resourceUrl"`
	CreatedAt   time.Time `gorm:"autoCreateTime"           json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"           json:"updatedAt"`
}

