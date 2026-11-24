package repositories

import (
	"wallet-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *WalletRepository) Credit(userID string, amount int64) (int64, error) {
	var wallet models.Wallet

	result := r.db.Clauses(clause.Returning{}).
		Model(&wallet).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return wallet.Balance, nil
}
