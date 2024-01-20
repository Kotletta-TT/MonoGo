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

// NewValidationError creates a new instance of ValidationError.
//
// Parameters:
//
//	get - an integer representing the HTTP status code for GET requests.
//	set - an integer representing the HTTP status code for SET requests.
//	err - a string describing the error message.
//
// Returns:
//
//	a pointer to a ValidationError struct.
func NewValidationErrror(get, set int, err string) *ValidationError {
	return &ValidationError{
		Err:           err,
		GetHTTPStatus: get,
		SetHTTPStatus: set,
	}
}

// Error returns the error message of the ValidationError.
//
// No parameters.
// string.
func (e *ValidationError) Error() string {
	return e.Err
}

// ValidateNameTypeParams validates the name and type parameters of the request.
//
// ctx: The gin context object.
// Returns a pointer to the Metrics struct and an error.
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

// ValidateValue validates the value parameter of a metric in a Gin context.
//
// Parameters:
// - ctx: The Gin context containing the HTTP request.
// - metric: The metric object to validate.
//
// Returns:
// - error: An error if the value is invalid, otherwise nil.
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

// ValidateParams validates the parameters of the given *gin.Context and returns a *common.Metrics object and an error.
//
// It calls the ValidateNameTypeParams function to validate the name and type parameters of the context.
// If there is an error during the validation, it returns nil and the error.
// If the request method is POST, it calls the ValidateValue function to validate the value parameter of the context.
// If there is an error during the validation, it returns nil and the error.
// Finally, it returns the validated *common.Metrics object and nil error.
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

// ValidateJSON validates the JSON data received in the request body and returns a Metrics object or an error.
//
// The function takes a gin.Context object as the parameter.
// It returns a pointer to a Metrics object and an error.
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

// ValidateBatchJSON validates and processes a batch of JSON metrics.
//
// ctx: The gin Context object.
// Returns: An array of metrics and an error.
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
