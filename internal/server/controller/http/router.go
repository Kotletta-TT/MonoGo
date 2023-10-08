package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func NewRouter(repo storage.Repository) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(RequestResponseLogging())
	engine.Use(gin.Recovery())
	engine.GET("/", ListMetrics(repo))
	engine.GET("/value/gauge/:metric", GetGaugeMetric(repo))
	engine.GET("/value/counter/:metric", GetCounterMetric(repo))
	engine.POST("/value/", GetJSONMetric(repo))
	engine.POST("/update/:metricType/:metric/:value", SetMetric(repo))
	engine.POST("/update/", SetJSONMetric(repo))
	return engine
}
