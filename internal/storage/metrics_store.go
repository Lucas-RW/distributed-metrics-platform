package storage

import (
	"sync"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/models"
)

var (
	mu    sync.Mutex
	store = make(map[string][]models.Metric)
)

func Save(metric models.Metric) error {
	mu.Lock()
	defer mu.Unlock()
	store[metric.Name] = append(store[metric.Name], metric)
	return nil
}

func GetAll(name string) []models.Metric {
	mu.Lock()
	defer mu.Unlock()
	return store[name]
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	store = make(map[string][]models.Metric)
}
