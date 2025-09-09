package repository

import (
	"kleio/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReleaseRepository struct {
	DB *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) *ReleaseRepository {
	return &ReleaseRepository{DB: db}
}

func (r *ReleaseRepository) CreateRelease(release *models.Release) error {
	return r.DB.Create(release).Error
}

func (r *ReleaseRepository) GetReleaseByID(id uuid.UUID) (*models.Release, error) {
	var release models.Release
	err := r.DB.Preload("Artists").
		Preload("Labels").
		Preload("Genres").
		Preload("Styles").
		Preload("Tracks").
		First(&release, "id = ?", id).Error
	return &release, err
}

func (r *ReleaseRepository) GetReleaseByDiscogsID(discogsID int) (*models.Release, error) {
	var release models.Release
	err := r.DB.Preload("Artists").
		Preload("Labels").
		Preload("Genres").
		Preload("Styles").
		Preload("Tracks").
		First(&release, "discogs_id = ?", discogsID).Error
	return &release, err
}

func (r *ReleaseRepository) UpdateRelease(release *models.Release) error {
	return r.DB.Save(release).Error
}

func (r *ReleaseRepository) DeleteRelease(id uuid.UUID) error {
	return r.DB.Delete(&models.Release{}, "id = ?", id).Error
}

func (r *ReleaseRepository) GetAllReleases(offset, limit int) ([]models.Release, error) {
	var releases []models.Release
	err := r.DB.Preload("Artists").
		Preload("Labels").
		Preload("Genres").
		Preload("Styles").
		Offset(offset).
		Limit(limit).
		Find(&releases).Error
	return releases, err
}

func (r *ReleaseRepository) SearchReleases(searchTerm string, offset, limit int) ([]models.Release, error) {
	var releases []models.Release
	searchPattern := "%" + searchTerm + "%"
	
	err := r.DB.Preload("Artists").
		Preload("Labels").
		Preload("Genres").
		Preload("Styles").
		Where("title ILIKE ?", searchPattern).
		Offset(offset).
		Limit(limit).
		Find(&releases).Error
	return releases, err
}

// Artist repository methods
func (r *ReleaseRepository) CreateArtist(artist *models.Artist) error {
	return r.DB.Create(artist).Error
}

func (r *ReleaseRepository) GetArtistByDiscogsID(discogsID int) (*models.Artist, error) {
	var artist models.Artist
	err := r.DB.First(&artist, "discogs_id = ?", discogsID).Error
	return &artist, err
}

func (r *ReleaseRepository) GetOrCreateArtist(artist *models.Artist) (*models.Artist, error) {
	existingArtist, err := r.GetArtistByDiscogsID(artist.DiscogsID)
	if err == nil {
		return existingArtist, nil
	}
	
	if err == gorm.ErrRecordNotFound {
		err := r.CreateArtist(artist)
		return artist, err
	}
	
	return nil, err
}

// Label repository methods
func (r *ReleaseRepository) CreateLabel(label *models.Label) error {
	return r.DB.Create(label).Error
}

func (r *ReleaseRepository) GetLabelByDiscogsID(discogsID int) (*models.Label, error) {
	var label models.Label
	err := r.DB.First(&label, "discogs_id = ?", discogsID).Error
	return &label, err
}

func (r *ReleaseRepository) GetOrCreateLabel(label *models.Label) (*models.Label, error) {
	existingLabel, err := r.GetLabelByDiscogsID(label.DiscogsID)
	if err == nil {
		return existingLabel, nil
	}
	
	if err == gorm.ErrRecordNotFound {
		err := r.CreateLabel(label)
		return label, err
	}
	
	return nil, err
}

// Genre repository methods
func (r *ReleaseRepository) GetOrCreateGenre(name string) (*models.Genre, error) {
	var genre models.Genre
	err := r.DB.Where("name = ?", name).First(&genre).Error
	
	if err == nil {
		return &genre, nil
	}
	
	if err == gorm.ErrRecordNotFound {
		genre = models.Genre{Name: name}
		err := r.DB.Create(&genre).Error
		return &genre, err
	}
	
	return nil, err
}

// Style repository methods
func (r *ReleaseRepository) GetOrCreateStyle(name string) (*models.Style, error) {
	var style models.Style
	err := r.DB.Where("name = ?", name).First(&style).Error
	
	if err == nil {
		return &style, nil
	}
	
	if err == gorm.ErrRecordNotFound {
		style = models.Style{Name: name}
		err := r.DB.Create(&style).Error
		return &style, err
	}
	
	return nil, err
}

// Track repository methods
func (r *ReleaseRepository) CreateTrack(track *models.Track) error {
	return r.DB.Create(track).Error
}

func (r *ReleaseRepository) GetTracksByReleaseID(releaseID uuid.UUID) ([]models.Track, error) {
	var tracks []models.Track
	err := r.DB.Where("release_id = ?", releaseID).Order("position").Find(&tracks).Error
	return tracks, err
}

func (r *ReleaseRepository) DeleteTracksByReleaseID(releaseID uuid.UUID) error {
	return r.DB.Where("release_id = ?", releaseID).Delete(&models.Track{}).Error
}