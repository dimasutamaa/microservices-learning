package services

import (
	"errors"
	"wallet-service/internal/dto"
	"wallet-service/internal/models"
	"wallet-service/internal/repositories"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo *repositories.WalletRepository
}

func NewWalletService(walletRepo *repositories.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (s *WalletService) CreateWallet(userID string) (*dto.WalletResponse, error) {
	wallet := &models.Wallet{
		ID:     uuid.NewString(),
		UserID: userID,
	}

	if _, err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return &dto.WalletResponse{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}, nil
}

func (s *WalletService) GetWalletBalance(userID string) (*dto.WalletResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.WalletResponse{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}, nil
}

func (s *WalletService) CreditWallet(userID string, amount int64) (int64, error) {
	newBalance, err := s.walletRepo.Credit(userID, amount)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (s *WalletService) Transfer(transactionID, fromUserID, toUserID string, amount int64) error {
	if fromUserID == toUserID {
		return ErrSameWallet
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.walletRepo.Transfer(transactionID, fromUserID, toUserID, amount)
}

var (
	ErrSameWallet = errors.New("sender and recipient must be different")
)
