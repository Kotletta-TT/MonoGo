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
	engine.Use(CompressMiddleware())
	engine.Use(gin.Recovery())
	engine.GET("/", ListMetrics(repo))
	engine.GET("/value/:metricType/:metric", GetMetric(repo))
	engine.POST("/value/", GetJSONMetric(repo))
	engine.POST("/update/:metricType/:metric/:value", SetMetric(repo))
	engine.POST("/update/", SetJSONMetric(repo))
	engine.POST("/updates/", SetBatchJSONMetric(repo))
	engine.GET("/ping", PingDB(repo))
	return engine
}
