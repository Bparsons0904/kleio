package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func getToken(r *http.Request) (string, error) {
	// Parse the request body
	var requestBody struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		return "", err
	}

	if requestBody.Token == "" {
		slog.Error("Token is required", "token", requestBody.Token)
		return "", nil
	}

	return requestBody.Token, nil
}
