// Package http implements some utils
package http

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// RunServer NewRouter creates a new gin.Engine instance and configures it with the provided repository and config.
//
// repo: The storage.Repository instance used for data storage.
// cfg: The config.Config instance used for configuration settings.
func RunServer(repo storage.Repository, cfg *config.Config) error {
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
	if cfg.SSL {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return err
		}
		caCert, err := os.ReadFile(cfg.CaPath)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool, ClientCAs: caCertPool, ClientAuth: tls.RequireAndVerifyClientCert}
		srv := http.Server{Handler: engine, TLSConfig: tlsCfg, Addr: cfg.RunServerAddr}
		return srv.ListenAndServeTLS("", "")
	} else {
		return engine.Run(cfg.RunServerAddr)
	}
}
