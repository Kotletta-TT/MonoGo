package common

import "testing"

func pointerInt64(i int64) *int64 {
	return &i
}

func TestSumDelta(t *testing.T) {
	type args struct {
		delta []*int64
	}
	tests := []struct {
		name string
		args args
		want *int64
	}{
		{
			"Normal",
			args{
				delta: []*int64{
					pointerInt64(1),
					pointerInt64(2),
					pointerInt64(3),
				},
			},
			pointerInt64(6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumDelta(tt.args.delta...); *got != *tt.want {
				t.Errorf("SumDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}
