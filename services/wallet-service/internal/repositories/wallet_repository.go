package repositories

import (
	"errors"
	"sort"
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

func (r *WalletRepository) Transfer(transactionID, fromUserID, toUserID string, amount int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.WalletTransfer
		err := tx.First(&existing, "transaction_id = ?", transactionID).Error
		if err == nil {
			if existing.FromUserID != fromUserID || existing.ToUserID != toUserID || existing.Amount != amount {
				return ErrTransferIDConflict
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		walletIDs := []string{fromUserID, toUserID}
		sort.Strings(walletIDs)
		var wallets []models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id IN ?", walletIDs).
			Order("user_id").
			Find(&wallets).Error; err != nil {
			return err
		}
		if len(wallets) != 2 {
			return ErrWalletNotFound
		}

		debit := tx.Model(&models.Wallet{}).
			Where("user_id = ? AND balance >= ?", fromUserID, amount).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount))
		if debit.Error != nil {
			return debit.Error
		}
		if debit.RowsAffected == 0 {
			var sender models.Wallet
			err := tx.First(&sender, "user_id = ?", fromUserID).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			if err != nil {
				return err
			}
			return ErrInsufficientFunds
		}

		credit := tx.Model(&models.Wallet{}).
			Where("user_id = ?", toUserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount))
		if credit.Error != nil {
			return credit.Error
		}
		if credit.RowsAffected == 0 {
			return ErrWalletNotFound
		}

		return tx.Create(&models.WalletTransfer{
			TransactionID: transactionID,
			FromUserID:    fromUserID,
			ToUserID:      toUserID,
			Amount:        amount,
		}).Error
	})
}

var (
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrTransferIDConflict = errors.New("transaction ID already used for a different transfer")
)
