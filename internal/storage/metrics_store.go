package storage

import (
    "context"
    "fmt"

    "github.com/Lucas-RW/distributed-metrics-platform/internal/models"
)

func Save(metric models.Metric) error {
    ctx := context.Background()

    _, err := pool.Exec(ctx,
        `INSERT INTO metrics (name, type, value, labels, timestamp)
         VALUES ($1, $2, $3, $4, $5)`,
        metric.Name,
        metric.Value,
        metric.Type,
        metric.Labels,
        metric.Timestamp,
    )
    if err != nil {
        return fmt.Errorf("failed to save metric: %w", err)
    }

    return nil
}

func GetAll(name string) ([]models.Metric, error) {
    ctx := context.Background()

    rows, err := pool.Query(ctx,
        `SELECT name, type, value, labels, timestamp
         FROM metrics
         WHERE name = $1
         ORDER BY timestamp DESC`,
        name,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to query metrics: %w", err)
    }
    defer rows.Close()

    var metrics []models.Metric
    for rows.Next() {
        var m models.Metric
        err := rows.Scan(
            &m.Name,
            &m.Value,
            &m.Type,
            &m.Labels,
            &m.Timestamp,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        metrics = append(metrics, m)
    }

    return metrics, nil
}