package services

import (
	"transaction-service/internal/client"
	"transaction-service/internal/models"
	"transaction-service/internal/repositories"

	"github.com/google/uuid"
)

type TransactionService struct {
	trxRepository *repositories.TransactionRepository
	walletClient  *client.WalletClient
}

func NewTransactionService(trxRepo *repositories.TransactionRepository, walletClient *client.WalletClient) *TransactionService {
	return &TransactionService{trxRepository: trxRepo, walletClient: walletClient}
}

func (s *TransactionService) Topup(userID string, amount int64) (any, error) {
	newBalance, err := s.walletClient.Credit(userID, amount)
	if err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		ID:     uuid.NewString(),
		UserID: userID,
		Amount: amount,
		Type:   "TOPUP",
		Status: "SUCCESS",
	}

	if err := s.trxRepository.Create(transaction); err != nil {
		return nil, err
	}

	return map[string]any{
		"transactionID": transaction.ID,
		"message":       "Top-up successful",
		"newBalance":    newBalance,
	}, nil
}
