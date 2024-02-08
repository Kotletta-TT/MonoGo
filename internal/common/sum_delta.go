// Package common implements some utils
package common

// SumDelta calculates the sum of the given delta values.
//
// The delta values are passed as variadic arguments of type *int64.
// The function returns a pointer to an int64, which is the sum of the delta values.
func SumDelta(delta ...*int64) *int64 {
	var sum int64
	for _, d := range delta {
		if d == nil {
			continue
		}
		sum += *d
	}
	return &sum
}
