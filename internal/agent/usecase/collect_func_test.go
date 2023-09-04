package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomMetricsCollector(t *testing.T) {

	tests := []struct {
		name           string
		wantCounter    int64
		wantLenMetrics int
	}{
		{
			name:           "One Iterate",
			wantCounter:    1,
			wantLenMetrics: 2,
		},
		{
			name:           "Two Iterate",
			wantCounter:    2,
			wantLenMetrics: 2,
		},
		{
			name:           "One Hungered Iterate",
			wantCounter:    100,
			wantLenMetrics: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.New()
			for i := 0; i < int(tt.wantCounter); i++ {
				CustomMetricsCollector(repo)
			}
			m := repo.GetMetrics()
			counter := m["PoolCount"].(int64)
			assert.Equal(t, tt.wantLenMetrics, len(m))
			assert.Equal(t, tt.wantCounter, counter)
		})
	}
}

func TestRuntimeMetricsCollector(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Regular Len",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.New()
			RuntimeMetricsCollector(repo)
			m := repo.GetMetrics()
			assert.Equal(t, 27, len(m))
		})
	}
}
