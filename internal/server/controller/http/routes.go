package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/Kotletta-TT/MonoGo/internal/server/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

func ListMetrics(repo storage.Repository) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		metrics := repo.GetAllMetrics()
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
