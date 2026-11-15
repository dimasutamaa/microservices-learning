package repositories

import (
	"wallet-service/internal/models"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) (*models.Wallet, error) {
	if err := r.db.Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *WalletRepository) FindByUserID(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.First(&wallet, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
