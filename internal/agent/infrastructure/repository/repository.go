package repository

import "sync"

type AgentRepository interface {
	StoreMetrics(map[string]interface{})
	GetMetrics() map[string]interface{}
}

type MemAgentRepository struct {
	mu      sync.Mutex
	storage map[string]interface{}
}

func New() AgentRepository {
	return &MemAgentRepository{
		mu:      sync.Mutex{},
		storage: make(map[string]interface{}),
	}
}

func (m *MemAgentRepository) StoreMetrics(metrics map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range metrics {
		switch v.(type) {
		case int64:
			oldVal, ok := m.storage[k]
			if ok {
				m.storage[k] = oldVal.(int64) + v.(int64)
			} else {
				m.storage[k] = v
			}
		case float64:
			m.storage[k] = v
		}
	}
}

func (m *MemAgentRepository) GetMetrics() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage
}
