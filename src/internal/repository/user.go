package repository

import (
	modelDb "bank_api/src/internal/models/db"
	"bank_api/src/pkg/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type User struct {
	database *db.DataBase
}

func NewUser(dataBase *db.DataBase) *User {
	return &User{
		dataBase,
	}
}

func (u *User) UpdateBalance(ctx context.Context, tx pgx.Tx, user modelDb.User) error {
	var queryUser = `UPDATE users SET balance = @Balance WHERE Id = @Id`
	var args = pgx.NamedArgs{
		"Id":      user.Id,
		"Balance": user.Balance,
	}

	if _, err := tx.Exec(ctx, queryUser, args); err != nil {
		return fmt.Errorf("ошибка обновления данных пользователя: %w", err)
	}

	return nil
}

func (u *User) GetBalance(ctx context.Context, tx pgx.Tx, userId modelDb.User) (float64, error) {
	var user modelDb.User
	var query = `SELECT balance FROM users WHERE id = @Id`
	var args = pgx.NamedArgs{"Id": userId.Id}

	if err := tx.QueryRow(ctx, query, args).Scan(&user.Balance); err != nil {
		return 0, fmt.Errorf("ошибка получения баланса пользователя: %w", err)
	}

	return user.Balance, nil
}
