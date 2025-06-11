package product

import (
	"context"
	"database/sql"
	"fmt"
	"frcofilippi/pedimeapp/shared/logger"
	"strings"

	"go.uber.org/zap"
)

func escapeLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func withUserContext(ctx context.Context, db *sql.DB, userId string, fn func(exec DBExecutor) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userLiteral := escapeLiteral(userId)
	setUserQuery := fmt.Sprintf("SET app.user_id = %s", userLiteral)

	_, err = tx.ExecContext(ctx, setUserQuery)
	if err != nil {
		tx.Rollback()
		logger.GetLogger().Error("database error", zap.String("message", err.Error()))
		return err
	}

	innerFunctionError := fn(tx)

	if innerFunctionError == sql.ErrNoRows {
		tx.Rollback()
		return innerFunctionError
	}

	resetUserStmt, err := tx.PrepareContext(ctx, "RESET app.user_id")
	if err != nil {
		tx.Rollback()
		logger.GetLogger().Error("database error", zap.String("message", err.Error()))
		return err
	}

	defer resetUserStmt.Close()

	_, resetErr := resetUserStmt.ExecContext(ctx)

	if resetErr != nil {
		tx.Rollback()
		logger.GetLogger().Error("database error", zap.String("message", resetErr.Error()))
		return resetErr
	}

	if innerFunctionError != nil {
		tx.Rollback()
		logger.GetLogger().Error("database error", zap.String("message", innerFunctionError.Error()))
		return innerFunctionError
	}

	return tx.Commit()
}

func (crp *WithUserProductRepository) GetById(ctx context.Context, id int64, userId string) (*Product, error) {
	var product *Product
	err := withUserContext(ctx, crp.db, userId, func(exec DBExecutor) error {
		var err error
		product, err = crp.innerRepo.GetById(ctx, exec, id, userId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (crp *WithUserProductRepository) Create(ctx context.Context, product *Product) (int64, error) {
	var id int64
	err := withUserContext(ctx, crp.db, product.UserId, func(exec DBExecutor) error {
		var err error
		id, err = crp.innerRepo.Create(ctx, exec, product)
		return err
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewProductRepositoryWithUser(db *sql.DB) (*WithUserProductRepository, error) {
	pgProductsRepo, err := NewProductRepository(db)
	if err != nil {
		return nil, err
	}
	return &WithUserProductRepository{
		innerRepo: pgProductsRepo,
		db:        db,
	}, nil
}
