package main

import (
	"fmt"
	"log"
	"net"

	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/rpc"
	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterUserServiceRPCServer(s, rpc.New())
	reflection.Register(s)

	address := "0.0.0.0:5050"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Backend Error %v", err)
	}

	fmt.Println("userService is listening on", address)

	db.ConnectDB()
	defer db.CloseDB()

	s.Serve(lis)
}
