package collectors

import (
	"testing"

	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/stretchr/testify/assert"
)

type MockCollectorStorage struct{}

func (m *MockCollectorStorage) StoreMetrics(map[string]*entity.Value) {}

func TestNewCollector(t *testing.T) {
	collector := NewCollector(&MockCollectorStorage{})
	assert.NotNil(t, collector)
}

func TestNewCollectorNilRepo(t *testing.T) {
	assert.Panics(t, assert.PanicTestFunc(func() {
		NewCollector(nil)
	}))
}

func TestCollector_RegisterCollectorMetricFunc(t *testing.T) {
	collector := NewCollector(&MockCollectorStorage{})
	collector.RegisterCollectorMetricFunc(func(repo collectorStorage) {})
	assert.Len(t, collector.registerCollectors, 1)
}

func TestCollector_Collect(t *testing.T) {
	collector := NewCollector(&MockCollectorStorage{})
	assert.NotPanics(t, collector.Collect)
}
