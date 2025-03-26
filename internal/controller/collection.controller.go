package controller

import (
	"log/slog"
)

func (c *Controller) SyncCollection() error {
	if err := c.SyncFolders(); err != nil {
		slog.Error("Failed to sync folders", "error", err)
		return err
	}

	if err := c.SyncReleases(); err != nil {
		slog.Error("Failed to sync collection", "error", err)
		return err
	}

	return nil
}

// func UpdateCollection(service database.Database) ([]Release, error) {
// 	db := service.GetDB()
// 	user, err := service.GetUser()
// 	if err != nil {
// 		slog.Error("Failed to get user", "error", err)
// 		return []Release{}, err
// 	}
//
// 	// folders, err := getDiscogFolders(user)
// 	// if err != nil {
// 	// 	slog.Error("Failed to get user folders", "error", err)
// 	// 	return []Release{}, err
// 	//
// 	// }
//
// 	lastSynced, err := getLocalFolderLastSynced(db)
// 	if err != nil {
// 		slog.Error("Failed to get last synced", "error", err)
// 	}
//
// 	now := time.Now().Add(-2 * time.Hour)
//
// 	if lastSynced.Before(now) {
// 		slog.Info("Last synced is older than 2 hours, updating folders...")
// 		updateFolders(db, folders)
// 	}
//
// 	for _, folder := range folders {
// 		updateCollectionByFolder(user, db, folder)
// 	}
//
// 	releases, err := GetAllReleases(db)
// 	if err != nil {
// 		slog.Error("Failed to get releases", "error", err)
// 		return []Release{}, err
// 	}
//
// 	return releases, nil
// }
