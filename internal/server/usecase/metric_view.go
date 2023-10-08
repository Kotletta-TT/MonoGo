package usecase

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
	"strconv"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

func TextPlainMetrics(metrics map[string]*shared.Metrics) []byte {
	textPlain := make([]byte, 0, 1024)
	for k, v := range metrics {
		switch v.MType {
		case GAUGE:
			stringValue := strconv.FormatFloat(*v.Value, 'f', -1, 64)
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %s\r\n", k, stringValue))...)
		case COUNTER:
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %d\r\n", k, *v.Delta))...)
		}
	}
	return textPlain
}

func TextPlainGaugeMetric(value float64) []byte {
	stringValue := strconv.FormatFloat(value, 'f', -1, 64)
	return []byte(fmt.Sprintf("%s\r\n", stringValue))
}

func TextPlainCounterMetrics(value int64) []byte {
	return []byte(fmt.Sprintf("%d\r\n", value))
}
