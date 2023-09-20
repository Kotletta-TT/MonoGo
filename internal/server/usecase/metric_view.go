package usecase

import (
	"fmt"
	"strconv"
)

func TextPlainMetrics(metrics map[string]interface{}) []byte {
	textPlain := make([]byte, 0, 1024)
	for k, v := range metrics {
		switch v := v.(type) {
		case float64:
			stringValue := strconv.FormatFloat(v, 'f', -1, 64)
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %s\r\n", k, stringValue))...)
		case int64:
			textPlain = append(textPlain, []byte(fmt.Sprintf("%s %d\r\n", k, v))...)
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
