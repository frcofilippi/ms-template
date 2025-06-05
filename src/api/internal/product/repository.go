package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
)

type PgProductRepository struct {
	db *sql.DB
}

func (pr *PgProductRepository) GetById(ctx context.Context, exec DBExecutor, id, customerId int64) (*Product, error) {
	query := `SELECT * FROM products p WHERE p.id = $1;`

	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()

	row := exec.QueryRowContext(ctx, query, id)

	var product Product

	err := row.Scan(&product.Id, &product.CustomerId, &product.Name, &product.Cost)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *PgProductRepository) Create(ctx context.Context, exec DBExecutor, product *Product) (int64, error) {
	query := "INSERT INTO products (customer_id, name, cost) VALUES ($1, $2, $3) RETURNING id"
	row := exec.QueryRowContext(ctx, query, product.CustomerId, product.Name, product.Cost)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	product.Id = id
	for _, event := range product.PendingEvents() {
		insertq := `INSERT INTO outbox_messages (aggregate_type, aggregate_id,event_type, payload)
						VALUES ($1,$2,$3,$4)`

		payload, _ := json.MarshalIndent(event, "", "  ")

		_, err := exec.ExecContext(ctx, insertq, "Product", product.Id, event.EventType(), payload)

		if err != nil {
			return 0, err
		}
	}

	product.CleanEvents()

	return id, nil
}

func NewProductRepository(db *sql.DB) (*PgProductRepository, error) {
	return &PgProductRepository{
		db: db,
	}, nil
}
