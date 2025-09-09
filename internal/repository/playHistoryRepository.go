package repository

import (
	"kleio/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayHistoryRepository struct {
	DB *gorm.DB
}

func NewPlayHistoryRepository(db *gorm.DB) *PlayHistoryRepository {
	return &PlayHistoryRepository{DB: db}
}

func (r *PlayHistoryRepository) CreatePlayHistory(play *models.PlayHistory) error {
	return r.DB.Create(play).Error
}

func (r *PlayHistoryRepository) GetPlayHistoryByID(id uuid.UUID) (*models.PlayHistory, error) {
	var play models.PlayHistory
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Stylus").
		First(&play, "id = ?", id).Error
	return &play, err
}

func (r *PlayHistoryRepository) GetUserPlayHistory(userID uuid.UUID, offset, limit int) ([]models.PlayHistory, error) {
	var plays []models.PlayHistory
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Stylus").
		Where("user_id = ?", userID).
		Order("played_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&plays).Error
	return plays, err
}

func (r *PlayHistoryRepository) GetUserPlayHistoryCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&models.PlayHistory{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *PlayHistoryRepository) GetReleasePlayHistory(userID, releaseID uuid.UUID) ([]models.PlayHistory, error) {
	var plays []models.PlayHistory
	err := r.DB.Preload("Stylus").
		Where("user_id = ? AND release_id = ?", userID, releaseID).
		Order("played_at DESC").
		Find(&plays).Error
	return plays, err
}

func (r *PlayHistoryRepository) GetRecentPlays(userID uuid.UUID, limit int) ([]models.PlayHistory, error) {
	var plays []models.PlayHistory
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Stylus").
		Where("user_id = ?", userID).
		Order("played_at DESC").
		Limit(limit).
		Find(&plays).Error
	return plays, err
}

func (r *PlayHistoryRepository) GetPlayHistoryByDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]models.PlayHistory, error) {
	var plays []models.PlayHistory
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Stylus").
		Where("user_id = ? AND played_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("played_at DESC").
		Find(&plays).Error
	return plays, err
}

func (r *PlayHistoryRepository) UpdatePlayHistory(play *models.PlayHistory) error {
	return r.DB.Save(play).Error
}

func (r *PlayHistoryRepository) DeletePlayHistory(id uuid.UUID) error {
	return r.DB.Delete(&models.PlayHistory{}, "id = ?", id).Error
}

func (r *PlayHistoryRepository) GetPlayCountByRelease(userID uuid.UUID) (map[uuid.UUID]int64, error) {
	type PlayCount struct {
		ReleaseID uuid.UUID `json:"release_id"`
		Count     int64     `json:"count"`
	}
	
	var results []PlayCount
	err := r.DB.Model(&models.PlayHistory{}).
		Select("release_id, count(*) as count").
		Where("user_id = ?", userID).
		Group("release_id").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	playCountMap := make(map[uuid.UUID]int64)
	for _, result := range results {
		playCountMap[result.ReleaseID] = result.Count
	}
	
	return playCountMap, nil
}

func (r *PlayHistoryRepository) GetMostPlayedReleases(userID uuid.UUID, limit int) ([]struct {
	Release   models.Release `json:"release"`
	PlayCount int64          `json:"play_count"`
}, error) {
	type Result struct {
		ReleaseID uuid.UUID `json:"release_id"`
		PlayCount int64     `json:"play_count"`
	}
	
	var results []Result
	err := r.DB.Model(&models.PlayHistory{}).
		Select("release_id, count(*) as play_count").
		Where("user_id = ?", userID).
		Group("release_id").
		Order("play_count DESC").
		Limit(limit).
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	var mostPlayed []struct {
		Release   models.Release `json:"release"`
		PlayCount int64          `json:"play_count"`
	}
	
	for _, result := range results {
		var release models.Release
		err := r.DB.Preload("Artists").First(&release, "id = ?", result.ReleaseID).Error
		if err != nil {
			continue
		}
		
		mostPlayed = append(mostPlayed, struct {
			Release   models.Release `json:"release"`
			PlayCount int64          `json:"play_count"`
		}{
			Release:   release,
			PlayCount: result.PlayCount,
		})
	}
	
	return mostPlayed, nil
}

func (r *PlayHistoryRepository) GetPlayStatsByMonth(userID uuid.UUID, year int) (map[int]int64, error) {
	type MonthlyStats struct {
		Month int   `json:"month"`
		Count int64 `json:"count"`
	}
	
	var results []MonthlyStats
	err := r.DB.Model(&models.PlayHistory{}).
		Select("EXTRACT(MONTH FROM played_at) as month, count(*) as count").
		Where("user_id = ? AND EXTRACT(YEAR FROM played_at) = ?", userID, year).
		Group("EXTRACT(MONTH FROM played_at)").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	monthlyStats := make(map[int]int64)
	for _, result := range results {
		monthlyStats[result.Month] = result.Count
	}
	
	return monthlyStats, nil
}