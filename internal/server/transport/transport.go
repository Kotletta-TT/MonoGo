package transport

import (
	"context"
	"fmt"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
)

const (
	GRPC = "grpc"
	HTTP = "http"
)

type Transport interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func NewTransport(cfg *config.Config, repo storage.Repository) (Transport, error) {
	switch cfg.Transport {
	case GRPC:
		return NewGRPCServer(cfg, repo)
	case HTTP:
		return NewHTTPServer(cfg, repo)
	}
	return nil, fmt.Errorf("invalid transport type: %s", cfg.Transport)
}
