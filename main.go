package main

import (
	"context"
	"log/slog"
	"os"
	"time"

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
	_, cancel := context.WithTimeout(context.Background(), 32*time.Second)
	defer cancel()
	var cfg pkg.Config

	if err := cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)

	}

}
