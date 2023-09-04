package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTTPSender_compileURL(t *testing.T) {
	type fields struct {
		repo   repository.AgentRepository
		client *resty.Client
	}
	type args struct {
		typeMetric string
		nameMetric string
		value      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Normal Counter Url",
			args:   args{value: "100", typeMetric: COUNTER, nameMetric: "counter1"},
			fields: fields{repo: repository.New(), client: resty.New()},
			want:   "http://localhost:8080/update/counter/counter1/100",
		},
		{
			name:   "Normal Gauge Url",
			args:   args{value: "100.0", typeMetric: GAUGE, nameMetric: "gauge1"},
			fields: fields{repo: repository.New(), client: resty.New()},
			want:   "http://localhost:8080/update/gauge/gauge1/100.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPSender{
				repo:       tt.fields.repo,
				client:     tt.fields.client,
				serverAddr: "localhost:8080",
			}
			assert.Equalf(t, tt.want, h.compileURL(tt.args.typeMetric, tt.args.nameMetric, tt.args.value), "compileURL(%v, %v, %v)", tt.args.typeMetric, tt.args.nameMetric, tt.args.value)
		})
	}
}
