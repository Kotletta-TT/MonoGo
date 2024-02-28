package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	v1 "github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
)

type HTTPServer struct {
	config *config.Config
	repo   storage.Repository
	router http.Handler
	h      *http.Server
}

func NewHTTPServer(config *config.Config, repo storage.Repository) (*HTTPServer, error) {
	return &HTTPServer{
		config: config,
		repo:   repo,
		router: v1.NewRouter(repo, config),
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context) error {
	s.h = &http.Server{
		Handler:      s.router,
		Addr:         s.config.RunServerAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if s.config.SSL {
		cert, err := tls.LoadX509KeyPair(s.config.CertPath, s.config.KeyPath)
		if err != nil {
			return err
		}
		caCert, err := os.ReadFile(s.config.CaPath)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool, ClientCAs: caCertPool, ClientAuth: tls.RequireAndVerifyClientCert}
		s.h.TLSConfig = tlsCfg
	}
	return s.h.ListenAndServeTLS("", "")
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.h.Shutdown(ctx)
}
