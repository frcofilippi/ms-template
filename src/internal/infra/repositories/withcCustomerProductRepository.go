package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"frcofilippi/pedimeapp/internal/business"
)

type WithCustomerProductRepository struct {
	innerRepo *PgProductRepository
	db        *sql.DB
}

type DBExecutor interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func withCustomerContext(ctx context.Context, db *sql.DB, customerId int64, fn func(exec DBExecutor) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, fmt.Sprintf("SET app.customer_id = %d", customerId))
	if err != nil {
		return err
	}

	innerFunctionError := fn(tx)

	if innerFunctionError == sql.ErrNoRows {
		tx.Rollback()
		return innerFunctionError
	}

	_, resetErr := tx.ExecContext(ctx, "RESET app.customer_id")
	if resetErr != nil && err == nil {
		err = resetErr
	}

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (crp *WithCustomerProductRepository) GetById(ctx context.Context, id, customerId int64) (*business.Product, error) {
	var product *business.Product
	err := withCustomerContext(ctx, crp.db, customerId, func(exec DBExecutor) error {
		var err error
		product, err = crp.innerRepo.GetById(ctx, exec, id, customerId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (crp *WithCustomerProductRepository) Create(ctx context.Context, product *business.Product) (int64, error) {
	var id int64
	err := withCustomerContext(ctx, crp.db, product.CustomerId, func(exec DBExecutor) error {
		var err error
		id, err = crp.innerRepo.Create(ctx, exec, product)
		return err
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewProductRepositoryWithCustomer(db *sql.DB) (*WithCustomerProductRepository, error) {
	pgProductsRepo, err := NewProductRepository(db)
	if err != nil {
		return nil, err
	}
	return &WithCustomerProductRepository{
		innerRepo: pgProductsRepo,
		db:        db,
	}, nil
}
