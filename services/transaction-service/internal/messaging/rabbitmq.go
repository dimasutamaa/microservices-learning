package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"transaction-service/internal/repositories"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	transferCommandQueue = "transaction.transfers"
	transferResultQueue  = "transaction.transfer-results"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	repository *repositories.TransactionRepository
}

func NewRabbitMQ(url string, repository *repositories.TransactionRepository) (*RabbitMQ, error) {
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
		connection: conn,
		channel:    ch,
		repository: repository,
	}

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

func (r *RabbitMQ) PublishTransfer(ctx context.Context, message TransferMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		"",
		transferCommandQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (r *RabbitMQ) ConsumeTransferResults() error {
	messages, err := r.channel.Consume(
		transferResultQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for message := range messages {
		var result TransferResult

		if err := json.Unmarshal(message.Body, &result); err != nil {
			log.Println("Invalid message:", err)
			message.Nack(false, false)
			continue
		}

		err := r.repository.UpdateStatus(result.TransactionID, result.Status, result.Error)
		if err != nil {
			log.Println("Failed to update transaction:", err)
			message.Nack(false, true)
			continue
		}

		message.Ack(false)
	}

	return fmt.Errorf("consumer stopped")
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.connection.Close()
}
