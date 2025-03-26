package controller

import "kleio/internal/database"

type Controller struct {
	DB database.Service
}

func InitNewController() *Controller {
	return &Controller{
		DB: database.New(),
	}
}
