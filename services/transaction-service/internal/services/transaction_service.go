package services

import (
	"context"
	"errors"
	"transaction-service/internal/client"
	"transaction-service/internal/dto"
	"transaction-service/internal/messaging"
	"transaction-service/internal/models"
	"transaction-service/internal/repositories"

	"github.com/google/uuid"
)

type TransactionService struct {
	trxRepository *repositories.TransactionRepository
	walletClient  *client.WalletClient
	rabbitMQ      *messaging.RabbitMQ
}

func NewTransactionService(trxRepo *repositories.TransactionRepository, walletClient *client.WalletClient, rabbitMQ *messaging.RabbitMQ) *TransactionService {
	return &TransactionService{trxRepository: trxRepo, walletClient: walletClient, rabbitMQ: rabbitMQ}
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

func (s *TransactionService) Transfer(ctx context.Context, fromUserID, toUserID string, amount int64) (*dto.TransferResponse, error) {
	if fromUserID == toUserID {
		return nil, ErrSelfTransfer
	}

	transaction := &models.Transaction{
		ID:       uuid.NewString(),
		UserID:   fromUserID,
		ToUserID: toUserID,
		Amount:   amount,
		Type:     "TRANSFER",
		Status:   "PENDING",
	}
	if err := s.trxRepository.Create(transaction); err != nil {
		return nil, err
	}

	message := messaging.TransferMessage{
		TransactionID: transaction.ID,
		FromUserID:    fromUserID,
		ToUserID:      toUserID,
		Amount:        amount,
	}
	if err := s.rabbitMQ.PublishTransfer(ctx, message); err != nil {
		if updateErr := s.trxRepository.UpdateStatus(transaction.ID, "FAILED", ErrTransferPublish.Error()); updateErr != nil {
			return nil, errors.Join(ErrTransferPublish, err, updateErr)
		}
		return nil, errors.Join(ErrTransferPublish, err)
	}

	return &dto.TransferResponse{
		TransactionID: transaction.ID,
		Status:        "PENDING",
		Message:       "Transfer initiated successfully.",
	}, nil
}

var (
	ErrSelfTransfer    = errors.New("sender and recipient must be different")
	ErrTransferPublish = errors.New("could not queue transfer")
)
