package migrations

import (
	"context"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Run(ctx context.Context, connString string) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed install dialect: %v", err)
	}

	if err := migrateStatus(ctx, connString); err != nil {
		return fmt.Errorf("failed get status migration: %v", err)
	}

	if err := migrateUp(ctx, connString); err != nil {
		return fmt.Errorf("failed migration up: %v", err)
	}

	return nil
}

func migrateUp(ctx context.Context, connString string) error {
	db, err := goose.OpenDBWithDriver("pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	return nil
}

func migrateDown(ctx context.Context, connString string) error {
	db, err := goose.OpenDBWithDriver("pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := goose.Down(db, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migrations: %v", err)
	}

	return nil
}

func migrateStatus(ctx context.Context, connString string) error {
	db, err := goose.OpenDBWithDriver("pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := goose.Status(db, "."); err != nil {
		return fmt.Errorf("failed to get migration status: %v", err)
	}

	return nil
}
