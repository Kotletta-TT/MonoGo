package common

import (
	"github.com/mailru/easyjson"
)

//go:generate easyjson -all internal/common/metrics.go

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetric(name string, mType string) *Metrics {
	return &Metrics{
		ID:    name,
		MType: mType,
	}
}

func (m *Metrics) String() string {
	buf, err := easyjson.Marshal(m)
	if err != nil {
		return ""
	}
	return string(buf)
}
