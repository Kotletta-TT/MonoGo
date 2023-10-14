package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/Kotletta-TT/MonoGo/internal/server/usecase"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

func ListMetrics(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		metrics, err := repo.GetAllMetrics()
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		byteMetrics := usecase.TextPlainMetrics(metrics)
		if _, err := ctx.Writer.Write(byteMetrics); err != nil {
			panic(err)
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func GetGaugeMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		mName := ctx.Param("metric")
		value, err := repo.GetGaugeMetric(mName)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)

			return
		}
		byteMetric := usecase.TextPlainGaugeMetric(value)
		if _, err = ctx.Writer.Write(byteMetric); err != nil {
			panic("Implement me")
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func GetCounterMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		mName := ctx.Param("metric")
		value, err := repo.GetCounterMetric(mName)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)

			return
		}
		byteMetric := usecase.TextPlainCounterMetrics(value)
		if _, err = ctx.Writer.Write(byteMetric); err != nil {
			panic("Implement me")
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func SetMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		switch ctx.Param("metricType") {
		case GAUGE:
			SetGaugeMetric(repo, ctx)
			return
		case COUNTER:
			SetCounterMetric(repo, ctx)
			return
		default:
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		}
	}
}

func SetGaugeMetric(repo storage.Repository, ctx *gin.Context) {
	name := ctx.Param("metric")
	mValue := ctx.Param("value")
	if mValue == "" {
		ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	value, err := usecase.ParseGaugeMetric(mValue)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	repo.StoreGaugeMetric(name, value)
	ctx.Writer.WriteHeader(http.StatusOK)
}

func SetCounterMetric(repo storage.Repository, ctx *gin.Context) {
	name := ctx.Param("metric")
	mValue := ctx.Param("value")
	if mValue == "" {
		ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	value, err := usecase.ParseCounterMetric(mValue)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	repo.StoreCounterMetric(name, value)
	ctx.Writer.WriteHeader(http.StatusOK)
}

func SetBatchJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		metrics := make([]*shared.Metrics, 0)
		batch := shared.SliceMetrics(metrics)
		err := easyjson.UnmarshalFromReader(ctx.Request.Body, &batch)
		if err != nil && err != io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = repo.StoreBatchMetric(batch)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}

func SetJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		m := shared.NewMetrics()
		err := easyjson.UnmarshalFromReader(ctx.Request.Body, m)
		if err != nil && err != io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		switch {
		case m.MType == GAUGE && m.Value != nil:
			repo.StoreGaugeMetric(m.ID, *m.Value)
			ctx.JSON(http.StatusOK, m)
		case m.MType == COUNTER && m.Delta != nil:
			repo.StoreCounterMetric(m.ID, *m.Delta)
			delta, err := repo.GetCounterMetric(m.ID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			m.Delta = &delta
			ctx.JSON(http.StatusOK, m)
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric type"})
		}
	}
}

func GetJSONMetric(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		m := shared.NewMetrics()
		err := easyjson.UnmarshalFromReader(ctx.Request.Body, m)
		if err != nil && err != io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		switch {
		case m.MType == GAUGE:
			val, err := repo.GetGaugeMetric(m.ID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			m.Value = &val
			ctx.JSON(http.StatusOK, m)
		case m.MType == COUNTER:
			delta, err := repo.GetCounterMetric(m.ID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			m.Delta = &delta
			ctx.JSON(http.StatusOK, m)
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric type"})
		}
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
