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
