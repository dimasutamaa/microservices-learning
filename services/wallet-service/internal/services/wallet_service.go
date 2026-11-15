package services

import (
	"wallet-service/internal/client"
	"wallet-service/internal/dto"
	"wallet-service/internal/models"
	"wallet-service/internal/repositories"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo *repositories.WalletRepository
	userClient *client.UserClient
}

func NewWalletService(walletRepo *repositories.WalletRepository, userClient *client.UserClient) *WalletService {
	return &WalletService{walletRepo: walletRepo, userClient: userClient}
}

func (s *WalletService) CreateWallet(userID string) (*dto.WalletResponse, error) {
	if _, err := s.userClient.FindUserByID(userID); err != nil {
		return nil, err
	}

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
