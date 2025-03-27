package controller

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

func GetUserIdentity(token string) (string, error) {
	url := fmt.Sprintf("%s/oauth/identity?token=%s", BaseURL, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "KleioApp/1.0 +https://github.com/bparsons0904/kleio")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"API returned non-200 status: %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var identity Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Found user identity: Username=%s, ID=%d", identity.Username, identity.ID)
	return identity.Username, nil
}
