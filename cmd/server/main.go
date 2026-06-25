package main

import (
    "context"
    "log"
    "net/http"
	"fmt"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/handlers"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/storage"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "os"
)

func main() {
    ctx := context.Background()

    if err := storage.InitDB(ctx); err != nil {
        log.Fatalf("could not connect to database: %v", err)
    }
    defer storage.CloseDB()

    if err := runMigrations(); err != nil {
        log.Fatalf("could not run migrations: %v", err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/ingest", handlers.IngestHandler)
    mux.Handle("/metrics", promhttp.Handler())

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}

func runMigrations() error {
    dbURL := os.Getenv("DATABASE_URL")

    m, err := migrate.New("file://migrations", dbURL)
    if err != nil {
        return fmt.Errorf("failed to create migrator: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    log.Println("Migrations ran successfully")
    return nil
}