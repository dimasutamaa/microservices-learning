package repositories

import (
	"transaction-service/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) UpdateStatus(trxID, status, msg string) error {
	result := r.db.Model(&models.Transaction{}).
		Where("id = ?", trxID).
		Updates(map[string]interface{}{
			"status":      status,
			"description": msg,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
