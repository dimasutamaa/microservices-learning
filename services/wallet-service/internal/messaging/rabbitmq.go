package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"wallet-service/internal/repositories"
	"wallet-service/internal/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	transferCommandQueue = "transaction.transfers"
	transferResultQueue  = "transaction.transfer-results"
)

type RabbitMQ struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	walletService *services.WalletService
}

func NewRabbitMQ(url string, walletService *services.WalletService) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	r := &RabbitMQ{
		connection:    conn,
		channel:       ch,
		walletService: walletService,
	}

	// Create queues
	_, err = ch.QueueDeclare(transferCommandQueue, true, false, false, false, nil)
	if err != nil {
		r.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(transferResultQueue, true, false, false, false, nil)
	if err != nil {
		r.Close()
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQ) PublishTransferResult(ctx context.Context, result TransferResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		"",
		transferResultQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    result.TransactionID,
			Body:         body,
		},
	)
}

func (r *RabbitMQ) ConsumeTransferCommands() error {
	// Process one message at a time
	err := r.channel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	messages, err := r.channel.Consume(
		transferCommandQueue,
		"",
		false, // Manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for message := range messages {
		var transfer TransferMessage

		if err := json.Unmarshal(message.Body, &transfer); err != nil {
			log.Println("Invalid transfer message:", err)
			message.Nack(false, false)
			continue
		}

		result, err := r.processTransfer(transfer)
		if err != nil {
			log.Println("Transfer processing failed:", err)
			message.Nack(false, true)
			continue
		}

		err = r.PublishTransferResult(context.Background(), result)
		if err != nil {
			log.Println("Failed to publish transfer result:", err)
			message.Nack(false, true)
			continue
		}

		message.Ack(false)
	}

	return fmt.Errorf("transfer consumer stopped")
}

func (r *RabbitMQ) processTransfer(message TransferMessage) (TransferResult, error) {
	result := TransferResult{
		TransactionID: message.TransactionID,
	}

	err := r.walletService.Transfer(
		message.TransactionID,
		message.FromUserID,
		message.ToUserID,
		message.Amount,
	)

	if err == nil {
		result.Status = "COMPLETED"
		return result, nil
	}

	if errors.Is(err, services.ErrSameWallet) ||
		errors.Is(err, repositories.ErrInsufficientFunds) ||
		errors.Is(err, repositories.ErrWalletNotFound) ||
		errors.Is(err, repositories.ErrTransferIDConflict) {

		result.Status = "FAILED"
		result.Error = err.Error()
		return result, nil
	}

	return TransferResult{}, err
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.connection.Close()
}
