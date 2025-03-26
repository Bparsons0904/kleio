package controller

import (
	"kleio/internal/database"
	"log/slog"
)

type Controller struct {
	DB   database.Database
	User database.User
}

func InitNewController() *Controller {
	return &Controller{
		DB: database.New(),
	}
}

func (c *Controller) SetUser() {
	user, err := c.DB.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return
	}

	c.User = user
}
