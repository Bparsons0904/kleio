package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.discogs.com"
)

type Folder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Count       int    `json:"count"`
	ResourceURL string `json:"resource_url"`
}

// FoldersResponse represents the response from the folders endpoint
type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

func QueryUserCollection(token string) {
	log.Println("Querying user collection...", "token", token)

	identity, err := GetUserIdentity(token)
	if err != nil {
		slog.Error("Failed to get user identity", "error", err)
		return
	}

	slog.Info("User identity", "identity", identity)

	// Build the URL for the folders endpoint
	url := fmt.Sprintf("%s/users/%s/collection/folders?token=%s", BaseURL, "deadstyle", token)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set required User-Agent header
	req.Header.Set("User-Agent", "KleioApp/1.0 +https://github.com/deadstyle/kleio")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		return
	}

	// Decode the response
	var foldersResp FoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&foldersResp); err != nil {
		log.Printf("Error decoding response: %v", err)
		return
	}

	// Log the folders
	log.Printf("Found %d folders", len(foldersResp.Folders))
	for _, folder := range foldersResp.Folders {
		log.Printf("Folder: ID=%d, Name=%s, Count=%d", folder.ID, folder.Name, folder.Count)
	}
}
