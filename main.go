package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/aws"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/database"
	"github.com/yaninyzwitty/sqs-postgres-microservice-inventory/internal/pkg"
)

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

	snsTopicArn, err := aws.CreateSnsTopicARN(ctx, "witty-topic", snsClient)
	if err != nil {
		slog.Error("Failed to create sns topic", "error", err)
		os.Exit(1)
	}
	fmt.Println(snsTopicArn)

	sqsClient, err := aws.LoadSQSClient(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load  and create sqs client", "error", err)
		os.Exit(1)
	}

	queueURL, err := aws.CreateQueueURL(ctx, "witty-queue", sqsClient)
	if err != nil {
		slog.Error("failed to create queue url", "error", err)
		os.Exit(1)
	}

	fmt.Println(queueURL)

}
