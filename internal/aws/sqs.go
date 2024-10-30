package aws

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func LoadSQSClient(ctx context.Context, region string) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load AWS sdk config", "error", err)
		return nil, err
	}
	sqsClient := sqs.NewFromConfig(cfg)
	return sqsClient, nil

}

func CreateQueueURL(ctx context.Context, queueName string, client *sqs.Client) (string, error) {
	createQueueOutput, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String("MyWittyQueue"),
	})
	if err != nil {
		log.Fatalf("unable to create SQS queue, %v", err)
	}
	queueURL := *createQueueOutput.QueueUrl
	fmt.Printf("Queue created, URL: %s\n", queueURL)
	return queueURL, nil

}
