package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemRepository_StoreGaugeMetric(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal Storage one metric",
			args: args{name: "gauge1", value: 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.StoreGaugeMetric(tt.args.name, tt.args.value)
			metric, _ := m.GetGaugeMetric(tt.args.name)
			assert.Equal(t, tt.args.value, metric)
		})
	}
}

func TestMemRepository_StoreCounterMetric(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal Storage one metric",
			args: args{name: "counter1", value: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.StoreCounterMetric(tt.args.name, tt.args.value)
			metric, _ := m.GetCounterMetric(tt.args.name)
			assert.Equal(t, tt.args.value, metric)
		})
	}
}
