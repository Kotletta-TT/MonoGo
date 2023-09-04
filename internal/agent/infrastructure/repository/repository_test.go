package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func MockMetrics(lenGauge, lenCounter int) map[string]interface{} {
	sum := lenGauge + lenCounter
	rndMetrics := make(map[string]interface{}, sum)
	for i := 0; i < sum; i++ {
		if lenGauge > 0 {
			rndMetrics[fmt.Sprintf("gauge%d", i)] = rand.Float64()
			lenGauge--
		}
		if lenCounter > 0 {
			rndMetrics[fmt.Sprintf("counter%d", i)] = int64(rand.Int())
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
				storage: make(map[string]interface{}),
			}
			genMetrics := MockMetrics(tt.lenGauge, tt.lenCounter)
			m.StoreMetrics(genMetrics)
			storeMetrics := m.GetMetrics()
			assert.Equal(t, len(storeMetrics), tt.lenCounter+tt.lenGauge)
			lc := tt.lenCounter
			lg := tt.lenGauge
			for _, v := range storeMetrics {
				switch v.(type) {
				case int64:
					tt.lenCounter--
					if tt.lenCounter < 0 {
						t.Errorf("len counter more than %d", lc)
					}
				case float64:
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
		waitCounter bool
		waitGauge   bool
		wantCounter int64
		wantGauge   float64
	}{
		{
			name:        "Normal one counter metrics 100",
			waitCounter: true,
			wantCounter: 100,
		},
		{
			name:        "Normal one gauge metrics 100.0",
			waitCounter: true,
			wantCounter: 100.0,
		},
		{
			name:        "Zero gauge metrics 0.0",
			waitCounter: true,
			wantCounter: 0.0,
		},
		{
			name:        "Zero counter metrics 0",
			waitCounter: true,
			wantCounter: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemAgentRepository{
				mu:      sync.Mutex{},
				storage: map[string]interface{}{"counter1": tt.wantCounter},
			}
			metrics := m.GetMetrics()
			if tt.waitGauge {
				floatVal := metrics["counter1"].(float64)
				assert.Equal(t, floatVal, tt.wantGauge)
			}
			if tt.waitCounter {
				intVal := metrics["counter1"].(int64)
				assert.Equal(t, intVal, tt.wantCounter)
			}
		})
	}
}
