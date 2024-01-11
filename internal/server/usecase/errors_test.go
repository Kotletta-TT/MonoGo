package usecase

import "testing"

func TestIncorrectValueMetrics_Error(t *testing.T) {
	type fields struct {
		Type  string
		Value string
		Err   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Normal",
			fields{
				Type:  "test",
				Value: "test",
				Err:   "test",
			},
			"testValue metric test incorrect from type test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := IncorrectValueMetrics{
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
				Err:   tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("IncorrectValueMetrics.Error() = <%s>, want <%s>", got, tt.want)
			}
		})
	}
}
