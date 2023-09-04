package entity

import (
	"log"
	"math"
)

type MetricKind int

const (
	// KindBad - undefined or incorrect metric type
	KindBad MetricKind = iota
	KindGauge
	KindCounter
)

type CustomMetricValue struct {
	value uint64
	kind  MetricKind
}

type CustomMetric struct {
	Name  string
	Value CustomMetricValue
}

func NewCustomGaugeMetric(name string, value float64) *CustomMetric {
	log.Printf("new gauge metric %s %f\n", name, value)
	return &CustomMetric{Name: name, Value: CustomMetricValue{kind: KindGauge, value: math.Float64bits(value)}}
}

func NewCustomCounterMetric(name string, value int64) *CustomMetric {
	log.Printf("new counter metric %s %d\n", name, value)
	return &CustomMetric{Name: name, Value: CustomMetricValue{kind: KindCounter, value: uint64(value)}}
}

func (m *CustomMetric) UpdateGageValue(value float64) {
	if m.Value.kind != KindGauge {
		panic("")
	}
	log.Printf("update gauge metric %s old: %f new: %f\n", m.Name, m.GetGaugeValue(), value)
	m.Value.value = math.Float64bits(value)

}

func (m *CustomMetric) UpdateCounterValue(value int64) {
	if m.Value.kind != KindCounter {
		panic("")
	}
	newCounter := int64(m.Value.value) + value
	log.Printf("update counter metric %s old: %d new: %d\n", m.Name, m.GetCounterValue(), newCounter)
	m.Value.value = uint64(newCounter)
}

func (m *CustomMetric) GetMetricKind() MetricKind {
	return m.Value.kind
}

func (m *CustomMetric) GetGaugeValue() float64 {
	if m.Value.kind != KindGauge {
		panic("")
	}
	return math.Float64frombits(m.Value.value)
}

func (m *CustomMetric) GetCounterValue() int64 {
	if m.Value.kind != KindCounter {
		panic("kind not counter")
	}
	return int64(m.Value.value)
}
