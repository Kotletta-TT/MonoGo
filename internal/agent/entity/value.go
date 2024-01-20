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

// NewValueFromUint64 creates a new Value object from a uint64 value and a KindMetric.
//
// Parameters:
// - v: the uint64 value to be converted into a Metric.
// - kind: the KindMetric to be assigned to the Value object.
//
// Return:
// - *Value: a pointer to the newly created Value object.
func NewValueFromUint64(v uint64, kind KindMetric) *Value {
	value := new(Value)
	value.Kind = kind
	convertFloat := float64(v)
	switch kind {
	case KindGauge:
		value.Metric = math.Float64bits(convertFloat)
	case KindCounter:
		value.Metric = v
	default:
		panic("unhandled default case")
	}
	return value
}

func NewValueFromUint32(v uint32, kind KindMetric) *Value {
	return NewValueFromUint64(uint64(v), kind)
}

// NewValueFromFloat64 creates a new Value object from a float64 value and a KindMetric.
//
// Parameters:
// - v: the float64 value to be converted into a Value object.
// - kind: the KindMetric indicating the type of the metric.
//
// Return:
// - value: the newly created Value object.
func NewValueFromFloat64(v float64, kind KindMetric) *Value {
	value := new(Value)
	value.Kind = kind
	switch kind {
	case KindGauge:
		value.Metric = math.Float64bits(v)
	case KindCounter:
		value.Metric = uint64(int64(v))
	default:
		panic("unhandled default case")
	}
	return value
}

// String returns a string representation of the Value.
//
// It converts the Value's Metric to a string based on its Kind:
// - For KindGauge, it converts the Metric to a float64 and returns it.
// - For KindCounter, it converts the Metric to an int64 and returns it.
//
// If the Kind is not KindGauge or KindCounter, it panics with
// "metric type unknown".
//
// Returns:
// - The string representation of the Value's Metric.
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
