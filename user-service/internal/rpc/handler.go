package rpc

import (
	"context"

	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/repo"
	services "github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/service"
	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/pb"
)

type server struct{}

// New - returns an instance of the UserServiceRPCServer
func New() pb.UserServiceRPCServer {
	return &server{}
}

func (s *server) GetService() services.UserService {
	userRepo := repo.New(db.DB)
	return services.New(userRepo)
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	svc := s.GetService()

	jwt, err := svc.LoginUser(ctx, req.Body.Name)
	if err != nil {
		return &pb.LoginResponse{}, err
	}

	return &pb.LoginResponse{
		Token: jwt,
	}, nil
}

func (s *server) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.GetUserBalanceResponse, error) {
	svc := s.GetService()

	balance, err := svc.GetUserBalance(ctx, req.Body.Token)

	if err != nil {
		return &pb.GetUserBalanceResponse{}, err
	}

	return &pb.GetUserBalanceResponse{
		Amount: balance,
	}, nil
}
