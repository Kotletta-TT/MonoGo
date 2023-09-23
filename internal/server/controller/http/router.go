package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func NewRouter(repo storage.Repository) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	//engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	//	return "test"
	//}))
	engine.Use(RequestResponseLogging())
	engine.Use(gin.Recovery())
	engine.GET("/", ListMetrics(repo))
	engine.GET("/value/gauge/:metric", GetGaugeMetric(repo))
	engine.GET("/value/counter/:metric", GetCounterMetric(repo))
	engine.POST("/update/:metricType/:metric/:value", SetMetric(repo))
	return engine
}
