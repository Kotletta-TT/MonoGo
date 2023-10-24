package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMetrics_UnmarshalEasyJSON(t *testing.T) {
	t.Run("UnmarshalEmptyArray", func(t *testing.T) {
		input := []byte("[]")
		ml := make(SliceMetrics, 0)
		ml.UnmarshalJSON(input)
		assert.Len(t, ml, 0)
	})

	t.Run("UnmarshalArrayWithOneMetric", func(t *testing.T) {
		input := []byte("[{\"id\": \"SomeID\", \"type\": \"gauge\", \"value\": 2}]")
		ml := make(SliceMetrics, 0)
		ml.UnmarshalJSON(input)
		assert.Len(t, ml, 1)
		assert.Equal(t, "SomeID", ml[0].ID)
		assert.Equal(t, "gauge", ml[0].MType)
		assert.Equal(t, float64(2), *ml[0].Value)
	})

	t.Run("UnmarshalArrayWithMultipleMetrics", func(t *testing.T) {
		input := []byte("[{\"id\": \"One\", \"type\": \"gauge\", \"value\": 1.0},{\"id\": \"Two\", \"type\": \"counter\", \"delta\": 2}, {\"id\": \"Three\", \"type\": \"gauge\", \"value\": 3.141592}]")
		ml := make(SliceMetrics, 0)
		ml.UnmarshalJSON(input)
		assert.Len(t, ml, 3)
		assert.Equal(t, "One", ml[0].ID)
		assert.Equal(t, "gauge", ml[0].MType)
		assert.Equal(t, float64(1.0), *ml[0].Value)
		assert.Equal(t, "Two", ml[1].ID)
		assert.Equal(t, "counter", ml[1].MType)
		assert.Equal(t, int64(2), *ml[1].Delta)
		assert.Equal(t, "Three", ml[2].ID)
		assert.Equal(t, "gauge", ml[2].MType)
		assert.Equal(t, float64(3.141592), *ml[2].Value)
	})
}

func TestSliceMetrics_MarshalJSON(t *testing.T) {
	ml := SliceMetrics{}
	expected1 := "[]"
	b, err := ml.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(b) != expected1 {
		t.Errorf("Expected %s, but got %s", expected1, string(b))
	}

	value := float64(1.1)
	ml = SliceMetrics{
		{
			ID:    "metric1",
			MType: "gauge",
			Value: &value,
		},
	}
	expected2 := `[{"id":"metric1","type":"gauge","value":1.1}]`
	b, err = ml.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(b) != expected2 {
		t.Errorf("Expected %s, but got %s", expected2, string(b))
	}

	value2 := 3.141592
	delta := int64(2)
	ml = SliceMetrics{
		{
			ID:    "metric1",
			MType: "gauge",
			Value: &value,
		},
		{
			ID:    "metric2",
			MType: "counter",
			Delta: &delta,
		},
		{
			ID:    "metric3",
			MType: "gauge",
			Value: &value2,
		},
	}
	expected3 := `[{"id":"metric1","type":"gauge","value":1.1},{"id":"metric2","type":"counter","delta":2},{"id":"metric3","type":"gauge","value":3.141592}]`
	b, err = ml.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(b) != expected3 {
		t.Errorf("Expected %s, but got %s", expected3, string(b))
	}
}
