package common

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(connectionStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")

	return db, nil
}
