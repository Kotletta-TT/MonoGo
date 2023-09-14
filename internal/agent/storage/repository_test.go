package storage

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"sync"
	"testing"
)

func MockMetrics(lenGauge, lenCounter int) map[string]*entity.Value {
	sum := lenGauge + lenCounter
	rndMetrics := make(map[string]*entity.Value, sum)
	for i := 0; i < sum; i++ {
		if lenGauge > 0 {

			rndMetrics[fmt.Sprintf("gauge%d", i)] = entity.NewValueFromFloat64(rand.Float64(), entity.KindGauge)
			lenGauge--
		}
		if lenCounter > 0 {
			rndMetrics[fmt.Sprintf("counter%d", i)] = entity.NewValueFromUint64(uint64(rand.Int()), entity.KindCounter)
			lenCounter--
		}
	}
	return rndMetrics
}

func TestMemAgentRepository_StoreMetrics(t *testing.T) {
	type args struct {
		metrics map[string]interface{}
	}
	tests := []struct {
		name       string
		lenGauge   int
		lenCounter int
	}{
		{
			name:       "Normal one gauge metrics",
			lenGauge:   1,
			lenCounter: 0,
		},
		{
			name:       "Normal one counter metrics",
			lenGauge:   0,
			lenCounter: 1,
		},
		{
			name:       "Normal one hungered metrics",
			lenGauge:   50,
			lenCounter: 50,
		},
		{
			name:       "Empty map add metrics",
			lenGauge:   0,
			lenCounter: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemAgentRepository{
				mu:      sync.Mutex{},
				storage: make(map[string]*entity.Value),
			}
			genMetrics := MockMetrics(tt.lenGauge, tt.lenCounter)
			m.StoreMetrics(genMetrics)
			storeMetrics := m.GetMetrics()
			assert.Equal(t, len(storeMetrics), tt.lenCounter+tt.lenGauge)
			lc := tt.lenCounter
			lg := tt.lenGauge
			for _, v := range storeMetrics {
				switch v.Kind {
				case entity.KindCounter:
					tt.lenCounter--
					if tt.lenCounter < 0 {
						t.Errorf("len counter more than %d", lc)
					}
				case entity.KindGauge:
					tt.lenGauge--
					if tt.lenGauge < 0 {
						t.Errorf("len gauge more than %d", lg)
					}
				}
			}
		})
	}
}

func TestMemAgentRepository_GetMetrics(t *testing.T) {
	tests := []struct {
		name        string
		metric      *entity.Value
		wantGauge   float64
		wantCounter int64
	}{
		{
			name:        "Normal one counter metrics 100",
			metric:      entity.NewValueFromUint64(100, entity.KindCounter),
			wantCounter: 100,
		},
		{
			name:      "Normal one gauge metrics 100.0",
			metric:    entity.NewValueFromFloat64(100.0, entity.KindGauge),
			wantGauge: 100.0,
		},
		{
			name:      "Zero gauge metrics 0.0",
			metric:    entity.NewValueFromFloat64(0.0, entity.KindGauge),
			wantGauge: 0.0,
		},
		{
			name:        "Zero counter metrics 0",
			metric:      entity.NewValueFromUint64(0, entity.KindCounter),
			wantCounter: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemAgentRepository{
				mu:      sync.Mutex{},
				storage: map[string]*entity.Value{"counter1": tt.metric},
			}
			metrics := m.GetMetrics()
			if tt.metric.Kind == entity.KindGauge {
				assert.Equal(t, math.Float64frombits(metrics["counter1"].Metric), tt.wantGauge)
			}
			if tt.metric.Kind == entity.KindCounter {
				assert.Equal(t, int64(metrics["counter1"].Metric), tt.wantCounter)
			}
		})
	}
}
