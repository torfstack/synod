package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pressly/goose/v3"
	"os"
)

var (
	DefaultMigrationsDir = "sql/migrations"
)

type MigrateOpts func(*MigrateOptions)

type MigrateOptions struct {
	MigrationsDir string
}

func WithMigrationsDir(dir string) MigrateOpts {
	return func(opts *MigrateOptions) {
		opts.MigrationsDir = dir
	}
}

func Migrate(ctx context.Context, connectionStr string, opts ...MigrateOpts) error {
	options := &MigrateOptions{
		MigrationsDir: DefaultMigrationsDir,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("could not set goose dialect 'postgres': %w", err)
	}

	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return fmt.Errorf("could not open database connection: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if _, err = os.Stat(options.MigrationsDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("specified migration directory '%s' does not exist", options.MigrationsDir)
	}

	if err = goose.UpContext(ctx, db, options.MigrationsDir); err != nil {
		return fmt.Errorf("could not run goose up: %w", err)
	}

	return nil
}
