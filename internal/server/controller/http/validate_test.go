package http

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/gin-gonic/gin"
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

func TestNewValidationErrror(t *testing.T) {
	type args struct {
		get int
		set int
		err string
	}
	tests := []struct {
		name string
		args args
		want *ValidationError
	}{
		{
			"Normal",
			args{
				get: 200,
				set: 200,
				err: "test",
			},
			&ValidationError{
				Err:           "test",
				GetHTTPStatus: 200,
				SetHTTPStatus: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValidationError(tt.args.get, tt.args.set, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidationErrror() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		Err           string
		GetHTTPStatus int
		SetHTTPStatus int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Normal",
			fields{
				Err:           "test",
				GetHTTPStatus: 200,
				SetHTTPStatus: 200,
			},
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationError{
				Err:           tt.fields.Err,
				GetHTTPStatus: tt.fields.GetHTTPStatus,
				SetHTTPStatus: tt.fields.SetHTTPStatus,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateNameTypeParams(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		{
			"Normal",
			args{
				ctx: &gin.Context{Params: []gin.Param{{Key: "metric", Value: "test"}, {Key: "metricType", Value: "gauge"}}},
			},
			&common.Metrics{ID: "test", MType: "gauge"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateNameTypeParams(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNameTypeParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateNameTypeParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateValue(t *testing.T) {
	type args struct {
		ctx    *gin.Context
		metric *common.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Normal",
			args{
				ctx:    &gin.Context{Params: []gin.Param{{Key: "metric", Value: "100"}, {Key: "metricType", Value: "gauge"}, {Key: "value", Value: "1.1"}}},
				metric: FabricGaugeMetric("100", 1.1),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateValue(tt.args.ctx, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("ValidateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateParams(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		{
			"Normal",
			args{
				ctx: &gin.Context{Request: &http.Request{Method: http.MethodPost}, Params: []gin.Param{{Key: "metric", Value: "100"}, {Key: "metricType", Value: "gauge"}, {Key: "value", Value: "1.1"}}},
			},
			FabricGaugeMetric("100", 1.1),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateParams(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		{
			"Normal",
			args{
				ctx: &gin.Context{Request: &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"id": "100", "type": "gauge", "value": 1.1}`)))}},
			},
			FabricGaugeMetric("100", 1.1),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJSON(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
