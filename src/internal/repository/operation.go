package repository

import (
	modelDb "bank_api/src/internal/models/db"
	"bank_api/src/pkg/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Operation struct {
	database *db.DataBase
}

func NewOperation(dataBase *db.DataBase) *Operation {
	return &Operation{
		dataBase,
	}
}

func (o *Operation) AddOperation(ctx context.Context, tx pgx.Tx, op modelDb.Operation) error {
	var queryOperation = `INSERT INTO operations (name, user_Id, scope) VALUES (@Name, @User_Id, @Scope)`
	var argsOperation = pgx.NamedArgs{
		"Name":    op.Name,
		"User_Id": op.UserId,
		"Scope":   op.Scope,
	}

	if _, err := tx.Exec(ctx, queryOperation, argsOperation); err != nil {
		return fmt.Errorf("ошибка добавления новой операции: %w", err)
	}

	return nil
}

func (o *Operation) GetLastOperations(ctx context.Context, u modelDb.User, count int) ([]modelDb.Operation, error) {
	var operations []modelDb.Operation
	var query = `SELECT * FROM operations WHERE user_Id = @User_Id ORDER BY date_create DESC LIMIT @Count`
	var args = pgx.NamedArgs{"User_Id": u.Id, "Count": count}

	rows, err := o.database.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения операций пользователя: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var op modelDb.Operation
		if err := rows.Scan(&op.Id, &op.Name, &op.UserId, &op.Scope, &op.DateCreate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования операции: %w", err)
		}
		operations = append(operations, op)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return operations, nil
}
