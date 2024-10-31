package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/aws"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/database"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/pkg"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/service"
)

var (
	queueURL    = "https://sqs.eu-north-1.amazonaws.com/651706749096/witty-queue"
	snsTopicArn = "arn:aws:sns:eu-north-1:651706749096:witty-topic"
)

// sqs_arn     = "arn:aws:sqs:eu-north-1:651706749096:witty-queue"

func main() {
	file, err := os.Open("config.yaml")
	if err != nil {
		slog.Error("failed to open config file", "error", err)
		os.Exit(1)
	}

	defer file.Close()
	if err != nil {
		slog.Error("failed to open config file", "error", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 32*time.Second)
	defer cancel()
	var cfg pkg.Config

	if err := cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)

	}

	db, err := database.NewDatabaseConnection(ctx, cfg.Database.DATABASE_URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	// ping the database
	err = database.PingDatabase(ctx, db)
	if err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	snsClient, err := aws.LoadSnsConfig(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load and create sns client", "error", err)
		os.Exit(1)
	}

	sqsClient, err := aws.LoadSQSClient(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load  and create sqs client", "error", err)
		os.Exit(1)
	}

	inventoryService := service.NewInventoryService(db, snsClient, sqsClient, &snsTopicArn, &queueURL)
	// process all order messages
	inventoryService.ProcessOrderMessage(ctx)

}
