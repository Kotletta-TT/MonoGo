package usecase

import (
	"testing"
)

func TestParseGaugeMetric(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Normal",
			args{
				value: "1.1",
			},
			1.1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGaugeMetric(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGaugeMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseGaugeMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCounterMetric(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"Normal",
			args{
				value: "1",
			},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCounterMetric(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCounterMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
