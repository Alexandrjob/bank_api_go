package service

import (
	modelDb "bank_api/src/internal/models/db"
	"bank_api/src/internal/models/dto"
	"bank_api/src/internal/repository"
	"bank_api/src/pkg/db"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type BankAccountService struct {
	dataBase  *db.DataBase
	user      *repository.User
	operation *repository.Operation
}

func NewBankAccountService(dataBase *db.DataBase, user *repository.User, operation *repository.Operation) *BankAccountService {
	return &BankAccountService{
		dataBase,
		user,
		operation,
	}
}

func (b *BankAccountService) Update(ctx context.Context, u dto.UpdateInfo) error {
	return b.dataBase.WithTransaction(ctx, func(tx pgx.Tx) error {
		var balance, err = b.user.GetBalance(ctx, tx, u.User)
		if err != nil {
			return err
		}

		u.User.Balance = balance + u.Operation.Scope
		if err := b.user.UpdateBalance(ctx, tx, u.User); err != nil {
			return err
		}

		u.Operation.DateCreate = time.Now()
		return b.operation.AddOperation(ctx, tx, u.Operation)
	})
}

func (b *BankAccountService) Transfer(ctx context.Context, t dto.TransferInfo) error {
	return b.dataBase.WithTransaction(ctx, func(tx pgx.Tx) error {
		err := b.updateSender(ctx, tx, t)
		if err != nil {
			return fmt.Errorf("ошибка обновления данных пользователя-отправителя: %w", err)
		}

		err = b.updateRecipient(ctx, tx, t)
		if err != nil {
			return fmt.Errorf("ошибка обновления данных пользователя-получателя: %w", err)
		}

		return nil
	})
}

func (b *BankAccountService) LastOperations(ctx context.Context, u modelDb.User, count int) ([]modelDb.Operation, error) {
	var result, err = b.operation.GetLastOperations(ctx, u, count)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BankAccountService) updateSender(ctx context.Context, tx pgx.Tx, t dto.TransferInfo) error {
	//Получение баланса пользователя-отправителя
	var user = modelDb.User{Id: t.SenderId}
	var balance, err = b.user.GetBalance(ctx, tx, user)
	if err != nil {
		return err
	}

	//Проверка, что на счету отправителя достаточно средств
	if balance-t.Scope < 0 {
		return errors.New("ошибка перевода средств: недостаточно денег на карте")
	}

	//Обновление счета отправителя
	user.Balance = balance - t.Scope
	if err := b.user.UpdateBalance(ctx, tx, user); err != nil {
		return err
	}

	//Добавление операции у отправителя
	var operation = modelDb.Operation{
		Name:       "Перевод",
		UserId:     t.SenderId,
		Scope:      t.Scope,
		DateCreate: time.Now(),
	}
	err = b.operation.AddOperation(ctx, tx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (b *BankAccountService) updateRecipient(ctx context.Context, tx pgx.Tx, t dto.TransferInfo) error {
	//Обновление счета получателя
	var user = modelDb.User{Id: t.RecipientId}
	var balance, err = b.user.GetBalance(ctx, tx, user)
	if err != nil {
		return err
	}

	user.Balance = balance + t.Scope
	if err := b.user.UpdateBalance(ctx, tx, user); err != nil {
		return err
	}

	//Добавление операции у получателя
	var operation = modelDb.Operation{
		Name:       "Получение",
		UserId:     t.RecipientId,
		Scope:      t.Scope,
		DateCreate: time.Now(),
	}
	err = b.operation.AddOperation(ctx, tx, operation)
	if err != nil {
		return err
	}

	return nil
}
