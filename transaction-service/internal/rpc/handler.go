package rpc

import (
	"context"

	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/repo"
	services "github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/service"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct{}

// New - returns an instance of the TransactionServiceRPCServer
func New() pb.TransactionServiceRPCServer {
	return &server{}
}

func (s *server) GetService(db *gorm.DB) services.TransactionService {
	txnRepo := repo.New(db)
	return services.New(txnRepo)
}

func (s *server) CreditAccount(ctx context.Context, req *pb.CreditAccountRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty
	db := db.DB
	svc := s.GetService(db)

	err := svc.CreditUserAccount(ctx, req.Body.Token, req.Body.Amount)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

func (s *server) DebitAccount(ctx context.Context, req *pb.DebitAccountRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty
	db := db.DB
	svc := s.GetService(db)

	err := svc.DebitUserAccount(ctx, req.Body.Token, req.Body.Amount)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
