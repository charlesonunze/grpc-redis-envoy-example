package rpc

import (
	"context"

	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/repo"
	services "github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/service"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct{}

// New - returns an instance of the TransactionServiceRPCServer
func New() pb.TransactionServiceRPCServer {
	return &server{}
}

func (s *server) GetService() services.TransactionService {
	userRepo := repo.New(db.DB)
	return services.New(userRepo)
}

func (s *server) CreditAccount(ctx context.Context, req *pb.CreditAccountRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	svc := s.GetService()
	err := svc.CreditUserAccount(ctx, req.Body.Token, req.Body.Amount)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

func (s *server) DebitAccount(ctx context.Context, req *pb.DebitAccountRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	svc := s.GetService()
	err := svc.DebitUserAccount(ctx, req.Body.Token, req.Body.Amount)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
