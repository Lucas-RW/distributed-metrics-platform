package storage

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDB(ctx context.Context) error {
    url := os.Getenv("DATABASE_URL")
    if url == "" {
        return fmt.Errorf("DATABASE_URL environment variable not set")
    }

    config, err := pgxpool.ParseConfig(url)
    if err != nil {
        return fmt.Errorf("failed to parse database URL: %w", err)
    }

    pool, err = pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return fmt.Errorf("failed to create connection pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    return nil
}

func CloseDB() {
    if pool != nil {
        pool.Close()
    }
}