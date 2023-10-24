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

func FabricGaugeMetric(name string, value float64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
}

func FabricCounterMetric(name string, delta int64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
}

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
