package handlers

import (
	"kleio/internal/database"
	"log/slog"
	"time"
)

const (
	BaseURL   = "https://api.discogs.com"
	UserAgent = "KleioApp/1.0 +https://github.com/bparsons0904/kleio"
)

func UpdateCollection(service database.Service) ([]Release, error) {
	db := service.GetDB()
	user, err := service.GetUser()
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return []Release{}, err
	}

	folders, err := GetFolders(user)
	if err != nil {
		slog.Error("Failed to get user folders", "error", err)
		return []Release{}, err

	}

	lastSynced, err := getLocalFolderLastSynced(db)
	if err != nil {
		slog.Error("Failed to get last synced", "error", err)
	}

	now := time.Now().Add(-2 * time.Hour)

	if lastSynced.Before(now) {
		slog.Info("Last synced is older than 2 hours, updating folders...")
		updateFolders(db, folders)
	}

	for _, folder := range folders {
		updateCollectionByFolder(user, db, folder)
	}

	releases, err := GetAllReleases(db)
	if err != nil {
		slog.Error("Failed to get releases", "error", err)
		return []Release{}, err
	}

	return releases, nil
}
