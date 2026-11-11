package repositories

import (
	"auth-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type AccessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) *AccessTokenRepository {
	return &AccessTokenRepository{db: db}
}

func (r *AccessTokenRepository) Create(token *models.AccessToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccessTokenRepository) FindByID(id string) (*models.AccessToken, error) {
	var token models.AccessToken
	if err := r.db.First(&token, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AccessTokenRepository) Revoke(id string, userID string) error {
	result := r.db.Model(&models.AccessToken{}).
		Where("id = ? AND user_id = ? AND revoked = ?", id, userID, false).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *AccessTokenRepository) IsUsable(id string, now time.Time) (*models.AccessToken, error) {
	var token models.AccessToken
	err := r.db.Where("id = ? AND revoked = ? AND expires_at > ?", id, false, now).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
