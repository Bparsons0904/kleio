package repository

import (
	"kleio/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StylusRepository struct {
	DB *gorm.DB
}

func NewStylusRepository(db *gorm.DB) *StylusRepository {
	return &StylusRepository{DB: db}
}

func (r *StylusRepository) CreateStylus(stylus *models.Stylus) error {
	return r.DB.Create(stylus).Error
}

func (r *StylusRepository) GetStylusByID(id uuid.UUID) (*models.Stylus, error) {
	var stylus models.Stylus
	err := r.DB.First(&stylus, "id = ?", id).Error
	return &stylus, err
}

func (r *StylusRepository) GetUserStyluses(userID uuid.UUID) ([]models.Stylus, error) {
	var styluses []models.Stylus
	err := r.DB.Where("user_id = ?", userID).
		Order("primary_stylus DESC, active DESC, name ASC").
		Find(&styluses).Error
	return styluses, err
}

func (r *StylusRepository) GetActiveStyluses(userID uuid.UUID) ([]models.Stylus, error) {
	var styluses []models.Stylus
	err := r.DB.Where("user_id = ? AND active = ?", userID, true).
		Order("primary_stylus DESC, name ASC").
		Find(&styluses).Error
	return styluses, err
}

func (r *StylusRepository) GetPrimaryStylus(userID uuid.UUID) (*models.Stylus, error) {
	var stylus models.Stylus
	err := r.DB.Where("user_id = ? AND primary_stylus = ?", userID, true).
		First(&stylus).Error
	return &stylus, err
}

func (r *StylusRepository) UpdateStylus(stylus *models.Stylus) error {
	return r.DB.Save(stylus).Error
}

func (r *StylusRepository) DeleteStylus(id uuid.UUID) error {
	return r.DB.Delete(&models.Stylus{}, "id = ?", id).Error
}

func (r *StylusRepository) SetPrimaryStylus(userID, stylusID uuid.UUID) error {
	// First, remove primary status from all styluses for this user
	err := r.DB.Model(&models.Stylus{}).
		Where("user_id = ?", userID).
		Update("primary_stylus", false).Error
	if err != nil {
		return err
	}
	
	// Then set the specified stylus as primary
	return r.DB.Model(&models.Stylus{}).
		Where("id = ? AND user_id = ?", stylusID, userID).
		Update("primary_stylus", true).Error
}

func (r *StylusRepository) GetStylusUsageStats(userID, stylusID uuid.UUID) (*StylusUsageStats, error) {
	type Result struct {
		PlayCount int64 `json:"play_count"`
	}
	
	var result Result
	err := r.DB.Model(&models.PlayHistory{}).
		Select("count(*) as play_count").
		Where("user_id = ? AND stylus_id = ?", userID, stylusID).
		Scan(&result).Error
	
	if err != nil {
		return nil, err
	}
	
	stylus, err := r.GetStylusByID(stylusID)
	if err != nil {
		return nil, err
	}
	
	stats := &StylusUsageStats{
		StylusID:         stylusID,
		Name:             stylus.Name,
		PlayCount:        result.PlayCount,
		ExpectedLifespan: stylus.ExpectedLifespan,
		UsagePercentage:  0,
	}
	
	if stylus.ExpectedLifespan > 0 {
		stats.UsagePercentage = float64(result.PlayCount) / float64(stylus.ExpectedLifespan) * 100
	}
	
	return stats, nil
}

func (r *StylusRepository) GetAllStylusUsageStats(userID uuid.UUID) ([]StylusUsageStats, error) {
	styluses, err := r.GetUserStyluses(userID)
	if err != nil {
		return nil, err
	}
	
	var stats []StylusUsageStats
	for _, stylus := range styluses {
		stylusStats, err := r.GetStylusUsageStats(userID, stylus.ID)
		if err != nil {
			continue // Skip this stylus if we can't get stats
		}
		stats = append(stats, *stylusStats)
	}
	
	return stats, nil
}

type StylusUsageStats struct {
	StylusID         uuid.UUID `json:"stylusId"`
	Name             string    `json:"name"`
	PlayCount        int64     `json:"playCount"`
	ExpectedLifespan int       `json:"expectedLifespan"`
	UsagePercentage  float64   `json:"usagePercentage"`
}