package repository

import (
	"kleio/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	DB *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{DB: db}
}

func (r *CollectionRepository) GetUserCollection(userID uuid.UUID, offset, limit int) ([]models.UserRelease, error) {
	var userReleases []models.UserRelease
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Release.Labels").
		Preload("Release.Genres").
		Preload("Release.Styles").
		Preload("Release.Tracks").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&userReleases).Error

	return userReleases, err
}

func (r *CollectionRepository) GetUserReleaseCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&models.UserRelease{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *CollectionRepository) GetUserReleaseByID(userID, releaseID uuid.UUID) (*models.UserRelease, error) {
	var userRelease models.UserRelease
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Release.Labels").
		Preload("Release.Genres").
		Preload("Release.Styles").
		Preload("Release.Tracks").
		Where("user_id = ? AND release_id = ?", userID, releaseID).
		First(&userRelease).Error
	
	return &userRelease, err
}

func (r *CollectionRepository) AddReleaseToCollection(userRelease *models.UserRelease) error {
	return r.DB.Create(userRelease).Error
}

func (r *CollectionRepository) RemoveReleaseFromCollection(userID, releaseID uuid.UUID) error {
	return r.DB.Where("user_id = ? AND release_id = ?", userID, releaseID).
		Delete(&models.UserRelease{}).Error
}

func (r *CollectionRepository) UpdateUserRelease(userRelease *models.UserRelease) error {
	return r.DB.Save(userRelease).Error
}

func (r *CollectionRepository) GetUserCollectionByFolder(userID uuid.UUID, folderID int, offset, limit int) ([]models.UserRelease, error) {
	var userReleases []models.UserRelease
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Release.Labels").
		Preload("Release.Genres").
		Preload("Release.Styles").
		Where("user_id = ? AND folder_id = ?", userID, folderID).
		Offset(offset).
		Limit(limit).
		Find(&userReleases).Error

	return userReleases, err
}

func (r *CollectionRepository) SearchUserCollection(userID uuid.UUID, searchTerm string, offset, limit int) ([]models.UserRelease, error) {
	var userReleases []models.UserRelease
	searchPattern := "%" + searchTerm + "%"
	
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Release.Labels").
		Joins("JOIN releases ON user_releases.release_id = releases.id").
		Joins("LEFT JOIN release_artists ra ON releases.id = ra.release_id").
		Joins("LEFT JOIN artists ON ra.artist_id = artists.id").
		Where("user_releases.user_id = ?", userID).
		Where("releases.title ILIKE ? OR artists.name ILIKE ?", searchPattern, searchPattern).
		Group("user_releases.user_id, user_releases.release_id").
		Offset(offset).
		Limit(limit).
		Find(&userReleases).Error

	return userReleases, err
}

func (r *CollectionRepository) GetUserReleasesByGenre(userID uuid.UUID, genreName string, offset, limit int) ([]models.UserRelease, error) {
	var userReleases []models.UserRelease
	
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Preload("Release.Genres").
		Joins("JOIN releases ON user_releases.release_id = releases.id").
		Joins("JOIN release_genres rg ON releases.id = rg.release_id").
		Joins("JOIN genres ON rg.genre_id = genres.id").
		Where("user_releases.user_id = ? AND genres.name = ?", userID, genreName).
		Offset(offset).
		Limit(limit).
		Find(&userReleases).Error

	return userReleases, err
}

func (r *CollectionRepository) GetUserReleasesByYear(userID uuid.UUID, year int, offset, limit int) ([]models.UserRelease, error) {
	var userReleases []models.UserRelease
	
	err := r.DB.Preload("Release").
		Preload("Release.Artists").
		Joins("JOIN releases ON user_releases.release_id = releases.id").
		Where("user_releases.user_id = ? AND releases.year = ?", userID, year).
		Offset(offset).
		Limit(limit).
		Find(&userReleases).Error

	return userReleases, err
}