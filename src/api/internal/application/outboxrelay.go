package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"frcofilippi/pedimeapp/shared/events"
	"log"
	"time"
)

type OutboxRelay struct {
	rmq       events.EventPublisher
	db        *sql.DB
	pollEvery time.Duration
}

func NewOutboxRelay(db *sql.DB, rmq events.EventPublisher) *OutboxRelay {
	return &OutboxRelay{
		db:        db,
		rmq:       rmq,
		pollEvery: time.Second * 2,
	}
}

func (or *OutboxRelay) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(or.pollEvery)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Default().Println("[OUTBOX] - Context canceled. Stoping outbox polling...")
			case <-ticker.C:
				err := or.PublishBatch(ctx)
				if err != nil {
					log.Default().Printf("[OUTBOX] - There was an error publishing messages")
				}
			}
		}
	}()
}

func (or *OutboxRelay) PublishBatch(ctx context.Context) error {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var messages []events.OutboxMessage

	query := `SELECT id, event_type, payload from outbox_messages where processed = false order by 1 asc`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	for rows.Next() {
		var msg events.OutboxMessage
		err := rows.Scan(&msg.Id, &msg.EventType, &msg.Payload)
		if err != nil {
			tx.Rollback()
			return err
		}
		messages = append(messages, msg)
	}

	for _, message := range messages {
		body, _ := json.Marshal(message)
		err = or.rmq.Publish(ctx, "", body)
		if err != nil {
			tx.Rollback()
			return err
		}

		updateq := `UPDATE outbox_messages SET processed = true, processed_at = $1 WHERE id = $2`
		_, err := tx.ExecContext(ctx, updateq, time.Now(), message.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
