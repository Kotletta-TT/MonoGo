package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/stretchr/testify/assert"
)

func FabricGaugeMetric(name string, value float64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
}

func FabricCounterMetric(name string, delta int64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
}

func MetricToInterface(m *common.Metrics) []interface{} {
	return []interface{}{m.ID, m.MType, m.Value, m.Delta}
}

func BatchMetricInterfaces(m []*common.Metrics) []interface{} {
	interfaceMetrics := make([]interface{}, 0, len(m))
	for _, metric := range m {
		interfaceMetrics = append(interfaceMetrics, MetricToInterface(metric)...)
	}
	return interfaceMetrics
}

func FabricBatchMetrics(gauge, counter int) []*common.Metrics {
	metrics := make([]*common.Metrics, 0, gauge+counter)
	for i := 0; i < gauge; i++ {
		metrics = append(metrics, FabricGaugeMetric(fmt.Sprintf("gauge%d", i), float64(i)))
	}
	for i := 0; i < counter; i++ {
		metrics = append(metrics, FabricCounterMetric(fmt.Sprintf("counter%d", i), int64(i)))
	}
	return metrics
}

func TestMemRepository_createOrUpdateMetric(t *testing.T) {
	testCases := []struct {
		name      string
		srcMetric *common.Metrics
		storage   map[string]*common.Metrics
		expMetric *common.Metrics
	}{
		{
			"create gauge metric",
			FabricGaugeMetric("gauge0", 1.1),
			make(map[string]*common.Metrics),
			FabricGaugeMetric("gauge0", 1.1),
		},
		{
			"update gauge metric",
			FabricGaugeMetric("gauge1", 2.2),
			map[string]*common.Metrics{
				"gauge1": FabricGaugeMetric("gauge1", 2.2),
			},
			FabricGaugeMetric("gauge1", 2.2),
		},
		{
			"create counter metric",
			FabricCounterMetric("counter0", 1),
			make(map[string]*common.Metrics),
			FabricCounterMetric("counter0", 1),
		},
		{
			"update counter metric",
			FabricCounterMetric("counter1", 2),
			map[string]*common.Metrics{
				"counter1": FabricCounterMetric("counter1", 2),
			},
			FabricCounterMetric("counter1", 4),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := MemRepository{mu: sync.Mutex{}, storage: tt.storage, cfg: &config.Config{StoreInterval: 1000}}
			m.createOrUpdateMetric(tt.srcMetric)
			assert.Equal(t, tt.expMetric, m.storage[tt.srcMetric.ID])
		})
	}
}

func TestMemRepository_StoreMetric(t *testing.T) {
	testCases := []struct {
		name    string
		src     *common.Metrics
		storage map[string]*common.Metrics
		exp     *common.Metrics
		wantErr bool
	}{
		{
			"Storage counter metric",
			FabricCounterMetric("counter0", 1),
			make(map[string]*common.Metrics),
			FabricCounterMetric("counter0", 1),
			false,
		},
		{
			"Storage gauge metric",
			FabricGaugeMetric("gauge0", 1.1),
			make(map[string]*common.Metrics),
			FabricGaugeMetric("gauge0", 1.1),
			false,
		},
		{
			"Update counter metric",
			FabricCounterMetric("counter1", 2),
			map[string]*common.Metrics{
				"counter1": FabricCounterMetric("counter1", 2),
			},
			FabricCounterMetric("counter1", 4),
			false,
		},
		{
			"Update gauge metric",
			FabricGaugeMetric("gauge1", 2),
			map[string]*common.Metrics{
				"gauge1": FabricGaugeMetric("gauge1", 2),
			},
			FabricGaugeMetric("gauge1", 2),
			false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := MemRepository{mu: sync.Mutex{}, storage: tt.storage, cfg: &config.Config{StoreInterval: 1000}}
			err := m.StoreMetric(tt.src)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.exp, m.storage[tt.src.ID])
			}
		})
	}
}

func TestMemRepository_GetMetric(t *testing.T) {
	testCases := []struct {
		name    string
		metric  *common.Metrics
		storage map[string]*common.Metrics
		exp     *common.Metrics
		wantErr bool
	}{
		{
			"Get gauge metric",
			FabricGaugeMetric("gauge0", 1.1),
			map[string]*common.Metrics{
				"gauge0": FabricGaugeMetric("gauge0", 1.1),
			},
			FabricGaugeMetric("gauge0", 1.1),
			false,
		},
		{
			"Get counter metric",
			FabricCounterMetric("counter0", 1),
			map[string]*common.Metrics{
				"counter0": FabricCounterMetric("counter0", 1),
			},
			FabricCounterMetric("counter0", 1),
			false,
		},
		{
			"Get non existing metric",
			FabricCounterMetric("counter1", 2),
			map[string]*common.Metrics{
				"counter0": FabricCounterMetric("counter0", 1),
			},
			nil,
			true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := MemRepository{mu: sync.Mutex{}, storage: tt.storage, cfg: &config.Config{StoreInterval: 1000}}
			err := m.GetMetric(tt.metric)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.exp, tt.metric)
			}
		})
	}
}

func TestMemRepository_GetListMetrics(t *testing.T) {
	testCases := []struct {
		name    string
		storage map[string]*common.Metrics
		exp     []*common.Metrics
	}{
		{
			"Get list metrics",
			map[string]*common.Metrics{
				"gauge0":   FabricGaugeMetric("gauge0", 1.1),
				"counter0": FabricCounterMetric("counter0", 1),
			},
			[]*common.Metrics{
				FabricGaugeMetric("gauge0", 1.1),
				FabricCounterMetric("counter0", 1),
			},
		},
		{
			"Get empty list metrics",
			make(map[string]*common.Metrics),
			make([]*common.Metrics, 0),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := MemRepository{mu: sync.Mutex{}, storage: tt.storage, cfg: &config.Config{StoreInterval: 1000}}
			storage, err := m.GetListMetrics()
			assert.NoError(t, err)
			assert.Equal(t, tt.exp, storage)
		})
	}
}
