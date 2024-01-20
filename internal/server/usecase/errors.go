package usecase

import (
	"fmt"
)

type IncorrectValueMetrics struct {
	Type  string
	Value string
	Err   string
}

// Error returns the error message for IncorrectValueMetrics.
//
// This function does not take any parameters.
// It returns a string representing the error message.
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

// Error returns the error message for IncorrectTypeMetrics.
//
// It concatenates the error message with the unsupported type.
// Returns the error message.
func (e IncorrectTypeMetrics) Error() string {
	e.Err += fmt.Sprintf("Type metric %s not supported", e.Type)
	return e.Err
}

type NoNameMetric struct{}

// Error returns the error message for NoNameMetric.
//
// No parameters.
// Returns a string.
func (e NoNameMetric) Error() string {
	return "Metric name not send"
}
