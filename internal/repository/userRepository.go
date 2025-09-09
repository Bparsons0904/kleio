package repository

import (
	"kleio/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "username = ?", username).Error
	return &user, err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) GetUserByToken(token string) (*models.User, error) {
	var authToken models.AuthToken
	err := r.DB.Preload("User").First(&authToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &authToken.User, nil
}

func (r *UserRepository) CreateAuthToken(token *models.AuthToken) error {
	return r.DB.Create(token).Error
}

func (r *UserRepository) GetAuthToken(token string) (*models.AuthToken, error) {
	var authToken models.AuthToken
	err := r.DB.Preload("User").First(&authToken, "token = ?", token).Error
	return &authToken, err
}

func (r *UserRepository) DeleteAuthToken(token string) error {
	return r.DB.Delete(&models.AuthToken{}, "token = ?", token).Error
}

func (r *UserRepository) GetUserAuthTokens(userID uuid.UUID) ([]models.AuthToken, error) {
	var tokens []models.AuthToken
	err := r.DB.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}