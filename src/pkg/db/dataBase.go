package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBase struct {
	Pool *pgxpool.Pool
}

func NewDataBase(pool *pgxpool.Pool) *DataBase {
	return &DataBase{
		pool,
	}
}

func (d *DataBase) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при создании транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
