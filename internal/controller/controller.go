package controller

import (
	"kleio/internal/database"
)

const (
	BaseURL   = "https://api.discogs.com"
	UserAgent = "KleioApp/1.0 +https://github.com/bparsons0904/kleio"
)

type Controller struct {
	DB database.Database
}

func InitNewController() *Controller {
	return &Controller{
		DB: database.New(),
	}
}
