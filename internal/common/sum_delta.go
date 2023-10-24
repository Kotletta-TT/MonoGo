package common

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
