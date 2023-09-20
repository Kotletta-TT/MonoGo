package usecase

import (
	"fmt"
)

type IncorrectValueMetrics struct {
	Type  string
	Value string
	Err   string
}

func (e IncorrectValueMetrics) Error() string {
	if e.Value != "" {
		e.Err += fmt.Sprintf("Value metric %v incorrect from type %s", e.Value, e.Type)
	}
	return e.Err
}

type IncorrectTypeMetrics struct {
	Type string
	Err  string
}

func (e IncorrectTypeMetrics) Error() string {
	e.Err += fmt.Sprintf("Type metric %s not supported", e.Type)
	return e.Err
}

type NoNameMetric struct{}

func (e NoNameMetric) Error() string {
	return "Metric name not send"
}
