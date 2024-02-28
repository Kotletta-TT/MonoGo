package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	pb "github.com/Kotletta-TT/MonoGo/internal/proto"
	grpcrouter "github.com/Kotletta-TT/MonoGo/internal/server/controller/grpc_router"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	config        *config.Config
	repo          storage.Repository
	rpcImplServer *grpcrouter.ProtoServer
	l             net.Listener
	g             *grpc.Server
}

func NewGRPCServer(config *config.Config, repo storage.Repository) (*GRPCServer, error) {
	return &GRPCServer{
		config:        config,
		repo:          repo,
		rpcImplServer: grpcrouter.NewProtoServer(repo),
	}, nil
}

func (s *GRPCServer) Start(ctx context.Context) error {
	var err error
	s.l, err = net.Listen("tcp", s.config.RunServerAddr)
	if err != nil {
		return err
	}
	s.g = grpc.NewServer()
	pb.RegisterMetricsServiceServer(s.g, s.rpcImplServer)
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
		s.l = tls.NewListener(s.l, tlsCfg)
	}
	return s.g.Serve(s.l)
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	s.g.Stop()
	return nil
}
