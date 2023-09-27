package sender

import (
	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/agent/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTTPSender_compileURL(t *testing.T) {
	type fields struct {
		repo   storage.AgentRepository
		client *resty.Client
	}
	type args struct {
		typeMetric entity.KindMetric
		nameMetric string
		value      uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Normal Counter Url",
			args:   args{value: 100, typeMetric: entity.KindCounter, nameMetric: "counter1"},
			fields: fields{repo: storage.New(), client: resty.New()},
			want:   "http://localhost:8080/update/counter/counter1/100",
		},
		{
			name:   "Normal Gauge Url",
			args:   args{value: 100, typeMetric: entity.KindGauge, nameMetric: "gauge1"},
			fields: fields{repo: storage.New(), client: resty.New()},
			want:   "http://localhost:8080/update/gauge/gauge1/100.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TextPlainSender{
				repo:   tt.fields.repo,
				client: tt.fields.client,
				cfg:    &config.Config{ServerHost: "localhost:8080"},
			}
			val := entity.NewValueFromUint64(tt.args.value, tt.args.typeMetric)
			assert.Equalf(t, tt.want, h.compileURL(tt.args.nameMetric, val), "compileURL(%v, %v)", tt.args.nameMetric, val)
		})
	}
}
