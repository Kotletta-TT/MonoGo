// Package usecase implements some utils
package usecase

import (
	"fmt"
	"strconv"

	"github.com/Kotletta-TT/MonoGo/internal/common"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

// FabricGaugeMetric generates a gauge metric with the given name and value.
//
// Parameters:
// - name: the ID of the metric
// - value: the value of the metric
//
// Returns:
// - *common.Metrics: a pointer to the generated metric
func FabricGaugeMetric(name string, value float64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
}

// FabricCounterMetric generates a Fabric counter metric.
//
// It takes in the name of the metric as a string and the delta value as an int64.
// It returns a pointer to a common.Metrics struct.
func FabricCounterMetric(name string, delta int64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
}

// TextPlainMetrics generates a text/plain representation of the given metrics.
//
// metrics is a slice of common.Metrics pointers.
// It returns a byte slice containing the text/plain representation of the metrics.
func TextPlainMetrics(metrics []*common.Metrics) []byte {
	textPlain := make([]byte, 0, 1024)
	for _, v := range metrics {
		switch v.MType {
		case GAUGE:
			stringValue := strconv.FormatFloat(*v.Value, 'f', -1, 64)
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %s\r\n", v.ID, stringValue))...)
		case COUNTER:
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %d\r\n", v.ID, *v.Delta))...)
		}
	}
	return textPlain
}

// TextPlainMetric generates a plain text metric based on the given Metrics object.
//
// metric: a pointer to a Metrics object containing the metric data.
// Returns: a byte array representing the generated plain text metric.
func TextPlainMetric(metric *common.Metrics) []byte {
	switch metric.MType {
	case GAUGE:
		stringValue := strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		return []byte(fmt.Sprintf("%s\r\n", stringValue))
	case COUNTER:
		return []byte(fmt.Sprintf("%d\r\n", *metric.Delta))
	}
	return nil
}
