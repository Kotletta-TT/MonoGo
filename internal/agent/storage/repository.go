package storage

import (
	"maps"
	"math"
	"sync"

	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/common"
)

type AgentRepository interface {
	StoreMetrics(map[string]*entity.Value)
	GetMetrics() map[string]*entity.Value
	GetMetricsSlice() []*common.Metrics
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
	return maps.Clone(m.storage)
}

func (m *MemAgentRepository) GetMetricsSlice() []*common.Metrics {
	oldMetrics := m.GetMetrics()
	metrics := make([]*common.Metrics, 0, len(oldMetrics))
	for k, v := range oldMetrics {
		newMetric := &common.Metrics{ID: k}
		switch v.Kind {
		case entity.KindCounter:
			newMetric.MType = "counter"
			delta := int64(v.Metric)
			newMetric.Delta = &delta
		case entity.KindGauge:
			newMetric.MType = "gauge"
			value := math.Float64frombits(v.Metric)
			newMetric.Value = &value
		}
		metrics = append(metrics, newMetric)
	}
	return metrics
}
