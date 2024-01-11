package entity

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewValueFromUint64_Gauge(t *testing.T) {
	value := NewValueFromUint64(100, KindGauge)
	assert.EqualValues(t, KindGauge, value.Kind)
	assert.Equal(t, float64(100), math.Float64frombits(value.Metric))
}

func TestNewValueFromUint64_Counter(t *testing.T) {
	value := NewValueFromUint64(100, KindCounter)
	assert.EqualValues(t, KindCounter, value.Kind)
	assert.Equal(t, 100, int(value.Metric))
}

func TestNewValueFromUint32_Gauge(t *testing.T) {
	value := NewValueFromUint32(100, KindGauge)
	assert.EqualValues(t, KindGauge, value.Kind)
	assert.Equal(t, float64(100), math.Float64frombits(value.Metric))
}

func TestNewValueFromUint32_Counter(t *testing.T) {
	value := NewValueFromUint32(100, KindCounter)
	assert.EqualValues(t, KindCounter, value.Kind)
	assert.Equal(t, 100, int(value.Metric))
}

func TestNewValueFromFloat64_Gauge(t *testing.T) {
	value := NewValueFromFloat64(100, KindGauge)
	assert.EqualValues(t, KindGauge, value.Kind)
	assert.Equal(t, float64(100), math.Float64frombits(value.Metric))
}

func TestNewValueFromFloat64_Counter(t *testing.T) {
	value := NewValueFromFloat64(100, KindCounter)
	assert.EqualValues(t, KindCounter, value.Kind)
	assert.Equal(t, 100, int(value.Metric))
}
