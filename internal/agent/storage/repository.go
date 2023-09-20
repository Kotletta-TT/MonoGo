package storage

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"sync"
)

type AgentRepository interface {
	StoreMetrics(map[string]*entity.Value)
	GetMetrics() map[string]*entity.Value
}

type MemAgentRepository struct {
	mu      sync.Mutex
	storage map[string]*entity.Value
}

func New() *MemAgentRepository {
	return &MemAgentRepository{
		mu:      sync.Mutex{},
		storage: make(map[string]*entity.Value),
	}
}

func (m *MemAgentRepository) StoreMetrics(metrics map[string]*entity.Value) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range metrics {
		switch v.Kind {
		case entity.KindCounter:
			oldVal, ok := m.storage[k]
			if ok {
				v.Metric = uint64(int64(oldVal.Metric) + int64(v.Metric))
				m.storage[k] = v
			} else {
				m.storage[k] = v
			}
		case entity.KindGauge:
			m.storage[k] = v
		default:
			panic("Parse bad metric Kind")
		}
	}
}

func (m *MemAgentRepository) GetMetrics() map[string]*entity.Value {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage
}
