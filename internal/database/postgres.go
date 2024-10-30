package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabaseConnection(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	if dbURL == "" {
		return nil, errors.New("please provide db url")
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("failed to create a connection pool", "error", err)
		return nil, err
	}

	return pool, nil
}

func PingDatabase(ctx context.Context, conn *pgxpool.Pool) error {
	const timeout = 60 * time.Second
	const pingInterval = 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		if err := conn.Ping(ctx); err != nil {
			slog.Error("failed to ping database", "error", err)
			time.Sleep(pingInterval)
			continue
		}
		slog.Info("pinged database succesfully")
		return nil
	}
	slog.Info("Completed 60 seconds of pinging the database without success")
	return fmt.Errorf("database not reachable within timeout period")

}
