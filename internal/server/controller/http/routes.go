// Package http implements some utils
package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/Kotletta-TT/MonoGo/internal/server/usecase"
	"github.com/gin-gonic/gin"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

// ListMetrics returns a function that handles the request to list metrics.
//
// The function takes a `storage.Repository` parameter `repo` which is used to retrieve the list of metrics.
// It returns a function that takes a `gin.Context` parameter `ctx` and writes the list of metrics to the response.
// The function also sets the content type header to "text/html".
// If an error occurs while retrieving the metrics or writing the response, the function sets the appropriate status code and returns.
// If the operation is successful, the function sets the status code to 200.
func ListMetrics(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		metrics, err := repo.GetListMetrics()
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		byteMetrics := usecase.TextPlainMetrics(metrics)
		if _, err := ctx.Writer.Write(byteMetrics); err != nil {
			ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

// GetMetric returns a Gin handler function that retrieves a metric from the specified repository.
//
// It takes a Gin context as input and validates the name and type parameters. If the validation fails,
// it returns a HTTP status code representing the validation error. If the validation succeeds, it
// calls the GetMetric method of the repository with the validated parameters. If the GetMetric
// operation fails, it returns a HTTP status code representing the failure. If the GetMetric
// operation succeeds, it converts the metric to a byte array of plain text format and writes it
// to the response writer of the context. Finally, it returns a HTTP status code representing a
// successful operation.
func GetMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		validatedGetMetric, err := ValidateNameTypeParams(ctx)
		if err != nil {
			var validateErr *ValidationError
			if errors.As(err, &validateErr) {
				ctx.Writer.WriteHeader(validateErr.GetHTTPStatus)
				return
			}
		}
		err = repo.GetMetric(validatedGetMetric)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		byteMetric := usecase.TextPlainMetric(validatedGetMetric)
		if _, err = ctx.Writer.Write(byteMetric); err != nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

// SetMetric is a function that returns a Gin handler function for setting a metric.
//
// It takes a repository of type storage.Repository as a parameter.
// The handler function receives a Gin context object as a parameter.
// It validates the parameters in the context using the ValidateParams function.
// If validation fails, it writes an appropriate HTTP status code to the response writer and returns.
// If validation succeeds, it stores the valid metric in the repository using the StoreMetric method.
// Finally, it writes an HTTP status code 200 (OK) to the response writer.
func SetMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		validMetric, err := ValidateParams(ctx)
		if err != nil {
			var validateErr *ValidationError
			if errors.As(err, &validateErr) {
				ctx.Writer.WriteHeader(validateErr.SetHTTPStatus)
				return
			}
		}
		repo.StoreMetric(validMetric)
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

// SetBatchJSONMetric returns a Gin handler function that sets a batch of JSON metrics in the repository.
//
// The function takes a Gin context `ctx` as its parameter.
// It validates the batch of JSON metrics using the `ValidateBatchJSON` function.
// If the validation fails, it returns a JSON response with the appropriate error status.
// If the validation succeeds, it prepares the batch for storage using the `PrepareBatchStore` function.
// It then stores the batch in the repository using the `StoreBatchMetric` method of the `repo` object.
// If the storage operation fails, it returns a JSON response with the appropriate error status.
func SetBatchJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		validateMetrics, err := ValidateBatchJSON(ctx)
		if err != nil {
			var validateErr *ValidationError
			if errors.As(err, &validateErr) {
				ctx.JSON(validateErr.SetHTTPStatus, gin.H{"error": err.Error()})
				return
			}
		}
		sortedBatch := usecase.PrepareBatchStore(validateMetrics)
		_, err = repo.StoreBatchMetric(sortedBatch)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}

// SetJSONMetric returns a function that handles a JSON metric and stores it in the repository.
//
// It takes a `repo` parameter of type `storage.Repository` which represents the repository to store the metric in.
// The returned function `ctx *gin.Context` is a handler function that accepts a `gin.Context` object.
// It validates the JSON metric received in the context, stores it in the repository, and returns an appropriate JSON response.
// The function returns no values.
func SetJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		validMetric, err := ValidateJSON(ctx)
		if err != nil {
			var validateErr *ValidationError
			if errors.As(err, &validateErr) {
				ctx.JSON(validateErr.SetHTTPStatus, gin.H{"error": err.Error()})
				return
			}
		}
		err = repo.StoreMetric(validMetric)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, validMetric)
	}
}

// GetJSONMetric returns a Gin handler function that retrieves a JSON metric.
//
// The function takes a `storage.Repository` as a parameter and returns a
// `func(ctx *gin.Context)` which is a handler function for Gin. The handler
// function retrieves a JSON metric from the context and validates it. If the
// validation fails, it returns an error response with the appropriate HTTP
// status code and error message. If the validation succeeds, it retrieves the
// metric from the repository and returns it as a JSON response with the HTTP
// status code 200 (OK).
func GetJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		validMetric, err := ValidateJSON(ctx)
		if err != nil {
			var validateErr *ValidationError
			if errors.As(err, &validateErr) {
				ctx.JSON(validateErr.GetHTTPStatus, gin.H{"error": err.Error()})
				return
			}
		}
		err = repo.GetMetric(validMetric)
		if err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, validMetric)
	}
}

// PingDB returns a function that pings the database and writes the appropriate status code to the response writer.
//
// The function takes a repository object as a parameter, which is used to perform the health check on the database.
// It returns a function that takes a Gin context object as a parameter.
// The function uses the context to create a new context with a timeout of 1 second and performs the health check using the repository object.
// If the health check fails, it writes a response with the status code 500 (Internal Server Error).
// If the health check succeeds, it writes a response with the status code 200 (OK).
func PingDB(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		pingCtx, cancel := context.WithTimeout(ctx.Request.Context(), 1*time.Second)
		defer cancel()
		err := repo.HealthCheck(pingCtx)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}
