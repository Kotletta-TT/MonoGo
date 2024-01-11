package common

import (
	"reflect"
	"testing"
)

func TestNewMetric(t *testing.T) {
	type args struct {
		name  string
		mType string
	}
	tests := []struct {
		name string
		args args
		want *Metrics
	}{
		{
			name: "Normal",
			args: args{
				name:  "test",
				mType: "gauge",
			},
			want: &Metrics{
				ID:    "test",
				MType: "gauge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetric(tt.args.name, tt.args.mType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_String(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Normal",
			fields: fields{
				ID:    "test",
				MType: "gauge",
				Delta: nil,
				Value: nil,
			},
			want: "{\"id\":\"test\",\"type\":\"gauge\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Metrics.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
