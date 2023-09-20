package entity

import (
	"fmt"
	"math"
)

const (
	KindBad = iota
	KindGauge
	KindCounter
)

type KindMetric int

type Value struct {
	Metric uint64
	Kind   KindMetric
}

func NewValueFromUint64(v uint64, kind KindMetric) *Value {
	value := new(Value)
	value.Kind = kind
	convertFloat := float64(v)
	switch kind {
	case KindGauge:
		value.Metric = math.Float64bits(convertFloat)
	case KindCounter:
		value.Metric = v
	}
	return value
}

func NewValueFromUint32(v uint32, kind KindMetric) *Value {
	return NewValueFromUint64(uint64(v), kind)
}

func NewValueFromFloat64(v float64, kind KindMetric) *Value {
	value := new(Value)
	value.Kind = kind
	switch kind {
	case KindGauge:
		value.Metric = math.Float64bits(v)
	case KindCounter:
		value.Metric = uint64(int64(v))
	}
	return value
}

func (v *Value) String() string {
	switch v.Kind {
	case KindGauge:
		return fmt.Sprintf("%f", math.Float64frombits(v.Metric))
	case KindCounter:
		return fmt.Sprintf("%d", int64(v.Metric))
	default:
		panic("metric type unknown")
	}
}
