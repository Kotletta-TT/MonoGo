package http

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/gin-gonic/gin"
)

// NewRouter creates a new gin.Engine instance and configures it with the provided repository and config.
//
// repo: The storage.Repository instance used for data storage.
// cfg: The config.Config instance used for configuration settings.
// Return: A pointer to the gin.Engine instance.
func NewRouter(repo storage.Repository, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(RequestResponseLogging())
	engine.Use(CompressMiddleware())
	engine.Use(HashSignMiddleWare(cfg))
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
