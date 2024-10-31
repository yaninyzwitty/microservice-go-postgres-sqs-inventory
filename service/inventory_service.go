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
	query := `INSERT INTO inventory (id, name, quantity, price) VALUES($1, $2, $3, $4)`

	res, err := s.db.Exec(ctx, query, order.ID, "witty-product-name", order.Quantity, 300.30)

	if err != nil {
		slog.Error("Failed to add inventory to db", "err", err)
		return err
	}

	if res.RowsAffected() == 0 {
		slog.Error("failed to insert row to db")
		return fmt.Errorf("failed to insert row to db")
	}

	// PUBLISH TO SNS FOR A NOTIFICATION
	message := fmt.Sprintf("Inventory created for item %s", order.ItemID)
	s.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: s.topicArn,
		Message:  &message,
	})

	return nil
}
