package collectors

import (
	"testing"

	"github.com/Kotletta-TT/MonoGo/internal/agent/storage"
	"github.com/stretchr/testify/assert"
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
			repo := storage.New()
			for i := 0; i < int(tt.wantCounter); i++ {
				CustomMetricsCollector(repo)
			}
			m := repo.GetMetrics()
			counter := m["PollCount"].Metric
			assert.Equal(t, tt.wantLenMetrics, len(m))
			assert.Equal(t, tt.wantCounter, int64(counter))
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
			repo := storage.New()
			RuntimeMetricsCollector(repo)
			m := repo.GetMetrics()
			assert.Equal(t, 27, len(m))
		})
	}
}
