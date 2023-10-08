package shared

import "github.com/mailru/easyjson"

//go:generate easyjson -all internal/shared/metrics.go

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) String() string {
	buf, err := easyjson.Marshal(m)
	if err != nil {
		return ""
	}
	return string(buf)
}
