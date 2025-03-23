package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Identity struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	ResourceURL string `json:"resource_url"`
}

// GetUserIdentity retrieves the identity of the user associated with the token
func GetUserIdentity(token string) (string, error) {
	// Build the URL for the identity endpoint
	url := fmt.Sprintf("%s/oauth/identity?token=%s", BaseURL, token)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set required User-Agent header
	req.Header.Set("User-Agent", "KleioApp/1.0 +https://github.com/yourusername/kleio")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"API returned non-200 status: %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	// Decode the response
	var identity Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Found user identity: Username=%s, ID=%d", identity.Username, identity.ID)
	return identity.Username, nil
}
