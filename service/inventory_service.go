package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/models"
)

type InventoryService struct {
	db        *pgxpool.Pool
	snsClient *sns.Client
	sqsClient *sqs.Client
	topicArn  *string
	queueUrl  *string
}

func NewInventoryService(db *pgxpool.Pool, snsClient *sns.Client, sqsClient *sqs.Client, topicArn *string, queueUrl *string) *InventoryService {
	return &InventoryService{
		db:        db,
		snsClient: snsClient,
		sqsClient: sqsClient,
		topicArn:  topicArn,
		queueUrl:  queueUrl,
	}
}

func (s *InventoryService) ProcessOrderMessage(ctx context.Context) error {
	for {
		output, err := s.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: s.queueUrl,
		})
		if err != nil {
			slog.Error("error receiving messages", "error", err)
			return err
		}
		for _, msg := range output.Messages {
			var order models.Order

			// json.Unmarshal([]byte(*msg.Body), &order)
			if err := json.Unmarshal([]byte(*msg.Body), &order); err != nil {
				slog.Error("failed to unmarshal data", "Error", err)
				return err
			}

			s.UpdateInventoryAndPublish(ctx, order)
			s.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      s.queueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			})

		}

	}
}

func (s *InventoryService) UpdateInventoryAndPublish(ctx context.Context, order models.Order) error {
	// insert into inventory
	var inventory models.Inventory
	query := `INSERT INTO inventory ( product_id, quantity ) VALUES($1, $2) RETURNING id, product_id, quantity, created_at`

	err := s.db.QueryRow(ctx, query, order.ProductId, order.Quantity).Scan(&inventory.Id, &inventory.ProductId, &inventory.Quantity, &inventory.CreatedAt)

	if err != nil {
		slog.Error("Failed to add inventory to db", "err", err)
		return fmt.Errorf("failed to insert inventory: %w", err)
	}

	// PUBLISH TO SNS FOR A NOTIFICATION
	message := fmt.Sprintf("Inventory created for item %s", order.ID)
	_, err = s.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: s.topicArn,
		Message:  &message,
	})
	if err != nil {
		slog.Error("Failed to publish SNS message", "order_id", order.ID, "err", err)
		return fmt.Errorf("failed to publish message to SNS: %w", err)

	}

	slog.Info("Successfully created inventory and published message", "product_id", order.ProductId, "quantity", order.Quantity)

	return nil
}
