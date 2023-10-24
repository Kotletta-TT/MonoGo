package http

import (
	"fmt"
	"net/http"

	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/Kotletta-TT/MonoGo/internal/server/usecase"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

type ValidationError struct {
	Err           string `json:"error"`
	GetHTTPStatus int
	SetHTTPStatus int
}

func NewValidationErrror(get, set int, err string) *ValidationError {
	return &ValidationError{
		Err:           err,
		GetHTTPStatus: get,
		SetHTTPStatus: set,
	}
}

func (e *ValidationError) Error() string {
	return e.Err
}

func ValidateNameTypeParams(ctx *gin.Context) (*common.Metrics, error) {
	name := ctx.Param("metric")
	mType := ctx.Param("metricType")
	if name == "" {
		return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric name is empty"))
	}
	if mType != GAUGE && mType != COUNTER {
		return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric type is not gauge or counter"))
	}
	return common.NewMetric(name, mType), nil
}

func ValidateValue(ctx *gin.Context, metric *common.Metrics) error {
	valueString := ctx.Param("value")
	switch metric.MType {
	case GAUGE:
		value, err := usecase.ParseGaugeMetric(valueString)
		if err != nil {
			return fmt.Errorf("%s %w", err, NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "invalid value for gauge"))
		}
		metric.Value = &value
	case COUNTER:
		value, err := usecase.ParseCounterMetric(valueString)
		if err != nil {
			return fmt.Errorf("%s %w", err, NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "invalid value for counter"))
		}
		metric.Delta = &value
	}
	return nil
}

func ValidateParams(ctx *gin.Context) (*common.Metrics, error) {
	metric, err := ValidateNameTypeParams(ctx)
	if err != nil {
		return nil, err
	}
	if ctx.Request.Method == http.MethodPost {
		err = ValidateValue(ctx, metric)
		if err != nil {
			return nil, err
		}
	}
	return metric, nil
}

func ValidateJSON(ctx *gin.Context) (*common.Metrics, error) {
	m := new(common.Metrics)
	err := easyjson.UnmarshalFromReader(ctx.Request.Body, m)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err, NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "invalid json"))
	}
	if m.ID == "" {
		return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric name is empty"))
	}
	if m.MType != GAUGE && m.MType != COUNTER {
		return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric type is not gauge or counter"))
	}
	return m, nil
}

func ValidateBatchJSON(ctx *gin.Context) ([]*common.Metrics, error) {
	metrics := make([]*common.Metrics, 0)
	batch := common.SliceMetrics(metrics)
	err := easyjson.UnmarshalFromReader(ctx.Request.Body, &batch)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err, NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "invalid json"))
	}
	for _, m := range batch {
		if m.ID == "" {
			return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric name is empty"))
		}
		if m.MType != GAUGE && m.MType != COUNTER {
			return nil, fmt.Errorf("%w", NewValidationErrror(http.StatusNotFound, http.StatusBadRequest, "metric type is not gauge or counter"))
		}
	}
	return batch, nil
}
