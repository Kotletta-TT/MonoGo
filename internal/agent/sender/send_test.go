package sender

import (
	"testing"

	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func FabricGaugeMetric(name string, value float64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
}

func FabricCounterMetric(name string, delta int64) *common.Metrics {
	return &common.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
}

func TestHTTPSender_compileURL(t *testing.T) {
	tests := []struct {
		name   string
		metric *common.Metrics
		want   string
	}{
		{
			name:   "Normal Counter Url",
			metric: FabricCounterMetric("counter1", 100),
			want:   "http://localhost:8080/update/counter/counter1/100",
		},
		{
			name:   "Normal Gauge Url",
			metric: FabricGaugeMetric("gauge1", 100.0),
			want:   "http://localhost:8080/update/gauge/gauge1/100.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TextPlainSender{
				repo:   nil,
				client: nil,
				cfg:    &config.Config{ServerHost: "localhost:8080"},
			}

			assert.Equalf(t, tt.want, h.compileURL(tt.metric), "compileURL(%v, %v)")
		})
	}
}

func TestNewRestyClient(t *testing.T) {
	cfg := config.NewConfig()
	client := NewRestyClient(cfg)
	assert.IsType(t, &resty.Client{}, client)
}

func TestNewHTTPSender_JSON(t *testing.T) {
	cnf := config.Config{SendType: JSON}
	snd, err := NewSender(nil, &cnf)
	assert.NoError(t, err)
	assert.IsType(t, &JSONSender{}, snd)

}

func TestNewHTTPSender_TextPlain(t *testing.T) {
	cnf := config.Config{SendType: TEXT}
	snd, err := NewSender(nil, &cnf)
	assert.NoError(t, err)
	assert.IsType(t, &TextPlainSender{}, snd)
}

func TestNewHTTPSender_Nil(t *testing.T) {
	cnf := config.Config{SendType: ""}
	assert.Panics(t, func() { NewSender(nil, &cnf) })
}
