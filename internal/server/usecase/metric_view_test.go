package usecase

import (
	"reflect"
	"testing"

	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestTextPlainView(t *testing.T) {
	testCases := []struct {
		desc    string
		metrics []*common.Metrics
		want    []byte
	}{
		{
			"One gauge metric",
			[]*common.Metrics{
				FabricGaugeMetric("gauge0", 1.1),
			},
			[]byte("gauge0 1.1\r\n"),
		},
		{
			"One counter metric",
			[]*common.Metrics{
				FabricCounterMetric("counter0", 1),
			},
			[]byte("counter0 1\r\n"),
		},
		{
			"Two metrics",
			[]*common.Metrics{
				FabricGaugeMetric("gauge0", 1.1),
				FabricCounterMetric("counter0", 1),
			},
			[]byte("gauge0 1.1\r\ncounter0 1\r\n"),
		},
		{
			"Empty metrics",
			[]*common.Metrics{},
			[]byte{},
		},
		{
			"Four metrics",
			[]*common.Metrics{
				FabricGaugeMetric("gauge0", 1.1),
				FabricGaugeMetric("gauge1", 2.2),
				FabricCounterMetric("counter0", 1),
				FabricCounterMetric("counter1", 2),
			},
			[]byte("gauge0 1.1\r\ngauge1 2.2\r\ncounter0 1\r\ncounter1 2\r\n"),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got := TextPlainMetrics(tt.metrics)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFabricGaugeMetric(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
		want *common.Metrics
	}{
		{
			"Normal",
			args{
				name:  "gauge0",
				value: 1.1,
			},
			FabricGaugeMetric("gauge0", 1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FabricGaugeMetric(tt.args.name, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FabricGaugeMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFabricCounterMetric(t *testing.T) {
	type args struct {
		name  string
		delta int64
	}
	tests := []struct {
		name string
		args args
		want *common.Metrics
	}{
		{
			"Normal",
			args{
				name:  "counter0",
				delta: 1,
			},
			FabricCounterMetric("counter0", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FabricCounterMetric(tt.args.name, tt.args.delta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FabricCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
