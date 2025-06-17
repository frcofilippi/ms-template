package common

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"frcofilippi/pedimeapp/shared/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewPostgresConnection(connectionStr string) (*sql.DB, error) {
	const (
		maxAttempts    = 5
		initialBackoff = 1 * time.Second
		maxBackoff     = 16 * time.Second
	)
	var (
		db      *sql.DB
		err     error
		backoff = initialBackoff
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("postgres", connectionStr)
		if err != nil {
			logger.GetLogger().Error(
				"error opening database connection",
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", maxAttempts),
				zap.String("error", err.Error()),
			)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = db.PingContext(ctx)
			if err == nil {
				fmt.Println("Connected to the database")
				return db, nil
			}
			logger.GetLogger().Error(
				"error pinging database",
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", maxAttempts),
				zap.String("error", err.Error()),
			)
			defer db.Close()
		}

		if attempt < maxAttempts {
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
	return nil, err
}
