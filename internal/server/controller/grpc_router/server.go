package grpcrouter

import (
	"context"

	"github.com/Kotletta-TT/MonoGo/internal/common"
	pb "github.com/Kotletta-TT/MonoGo/internal/proto"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProtoServer struct {
	repo storage.Repository
	pb.UnimplementedMetricsServiceServer
}

func (s *ProtoServer) GetMetric(ctx context.Context, receiveMetric *pb.GetMetricRequest) (*pb.Metric, error) {
	metric := &common.Metrics{}
	metric.ID = receiveMetric.Id
	metric.MType = receiveMetric.Mtype.String()
	err := s.repo.GetMetric(metric)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, `Metric %s not found`, receiveMetric.Id)
	}
	return metric.ToProto(), nil
}

func (s *ProtoServer) SetMetric(ctx context.Context, receiveMetric *pb.Metric) (*pb.Metric, error) {
	metric := common.NewMetricFromProto(receiveMetric)
	err := s.repo.StoreMetric(metric)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `Metric %s not stored`, receiveMetric.Name)
	}
	return metric.ToProto(), nil
}

func (s *ProtoServer) SetBulkMetrics(ctx context.Context, receiveMetric *pb.SetBulkMetricsRequest) (*pb.Empty, error) {
	metrics := common.NewSliceMetricsFromProto(receiveMetric)
	_, err := s.repo.StoreBatchMetric(metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `Metrics %s not stored`, receiveMetric.Metrics)
	}
	return &pb.Empty{}, nil
}

func (s *ProtoServer) GetListMetrics(ctx context.Context, receiveMetric *pb.Empty) (*pb.GetListMetricsResponse, error) {
	metrics, err := s.repo.GetListMetrics()
	if err != nil {
		return nil, status.Errorf(codes.Internal, `Metrics %s not stored`, receiveMetric)
	}
	pbMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, m := range metrics {
		pbMetrics = append(pbMetrics, m.ToProto())
	}
	return &pb.GetListMetricsResponse{Metrics: pbMetrics}, nil
}

func NewProtoServer(repo storage.Repository) *ProtoServer {
	return &ProtoServer{repo: repo}
}
