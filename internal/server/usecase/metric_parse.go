package usecase

import (
	"strconv"
)

// ParseGaugeMetric parses a string value into a float64.
//
// It takes in a string value and returns a float64 and an error.
func ParseGaugeMetric(value string) (float64, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return floatValue, nil
}

// ParseCounterMetric parses a string value and returns its equivalent int64 representation.
//
// It takes a string value as input and attempts to parse it using the strconv.ParseInt function.
// If the parsing is successful, it returns the parsed int64 value.
// If the parsing fails, it returns an error indicating the reason for the failure.
//
// Parameters:
// - value: The string value to be parsed.
//
// Returns:
// - int64: The parsed int64 value.
// - error: An error indicating the reason for the parsing failure.
func ParseCounterMetric(value string) (int64, error) {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
