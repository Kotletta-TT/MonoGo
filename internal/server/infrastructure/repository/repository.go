package repository

import (
	"fmt"
	"sync"
)

type Repository interface {
	StoreGaugeMetric(name string, value float64)
	StoreCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	GetAllMetrics() map[string]interface{}
}

type MemRepository struct {
	mu      sync.Mutex
	storage map[string]interface{}
}

func New() Repository {
	return &MemRepository{mu: sync.Mutex{}, storage: make(map[string]interface{}, 30)}
}

func (m *MemRepository) StoreGaugeMetric(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[name] = value
}

func (m *MemRepository) StoreCounterMetric(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.storage[name]
	if !ok {
		m.storage[name] = value
		return
	}
	valInt, ok := val.(int64)
	if ok {
		m.storage[name] = value + valInt
	}
}

func (m *MemRepository) GetGaugeMetric(name string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.storage[name]
	if !ok {
		return 0.0, fmt.Errorf("metric not found")
	}
	val, ok := metric.(float64)
	if !ok {
		return 0.0, fmt.Errorf("metric format error")
	}
	return val, nil
}

func (m *MemRepository) GetCounterMetric(name string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.storage[name]
	if !ok {
		return 0.0, fmt.Errorf("metric not found")
	}
	val, ok := metric.(int64)
	if !ok {
		return 0.0, fmt.Errorf("metric format error")
	}
	return val, nil
}

func (m *MemRepository) GetAllMetrics() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage
}
