package controller

import (
	"log/slog"
)

func (c *Controller) GetAuth() (payload Payload, err error) {
	payload.Token, err = c.DB.GetToken()
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return Payload{}, nil
	}

	err = payload.GetPayload(c)
	if err != nil {
		slog.Error("Failed to get payload", "error", err)
		return payload, err
	}

	return payload, nil
}

func (c *Controller) SaveToken(token string) (payload Payload, err error) {
	username, err := GetUserIdentity(token)
	if err != nil {
		slog.Error("Failed to get user identity", "error", err)
		return payload, err
	}

	err = c.DB.SaveToken(token, username)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		return payload, err
	}

	payload, err = c.GetAuth()
	if err != nil {
		slog.Error("Failed to get auth", "error", err)
		return payload, err
	}

	go c.asyncCollection()
	payload.SyncingData = true

	return payload, nil
}
