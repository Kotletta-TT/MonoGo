// Package common implements some utils
package common

import (
	pb "github.com/Kotletta-TT/MonoGo/internal/proto"
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

// NewMetric initializes and returns a new Metrics object.
//
// Parameters:
//   - name: the name of the metric
//   - mType: the type of the metric
//
// Returns:
//   - *Metrics: a pointer to the newly created Metrics object
func NewMetric(name string, mType string) *Metrics {
	return &Metrics{
		ID:    name,
		MType: mType,
	}
}

func NewMetricFromProto(receiveMetric *pb.Metric) *Metrics {
	m := NewMetric(receiveMetric.Name, receiveMetric.Mtype.String())
	if receiveMetric.Delta != nil {
		m.Delta = receiveMetric.Delta
	}
	if receiveMetric.Value != nil {
		m.Value = receiveMetric.Value
	}
	return m
}

func NewSliceMetricsFromProto(receiveMetric *pb.SetBulkMetricsRequest) []*Metrics {
	metrics := make([]*Metrics, 0, len(receiveMetric.Metrics))
	for _, metric := range receiveMetric.Metrics {
		metrics = append(metrics, NewMetricFromProto(metric))
	}
	return metrics
}

func NewProtoMetricsFromSlice(metrics []*Metrics) *pb.SetBulkMetricsRequest {
	pbMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, metric := range metrics {
		pbMetrics = append(pbMetrics, metric.ToProto())
	}
	return &pb.SetBulkMetricsRequest{Metrics: pbMetrics}
}

// String returns the string representation of the Metrics struct.
//
// It marshals the Metrics struct using the easyjson.Marshal function.
// If an error occurs during marshaling, it returns an empty string.
// Otherwise, it returns the marshaled struct as a string.
func (v *Metrics) String() string {
	buf, err := easyjson.Marshal(v)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (v *Metrics) ToProto() *pb.Metric {
	return &pb.Metric{
		Name:  v.ID,
		Mtype: pb.MType(pb.MType_value[v.MType]),
		Delta: v.Delta,
		Value: v.Value,
	}
}