package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/repo"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/rpc"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/service"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/pb"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PWD"),
	})
	GET_USER_BALANCE_TOPIC = "get_user_balance"
	USER_BALANCE_TOPIC     = "user_balance"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterTransactionServiceRPCServer(s, rpc.New())
	reflection.Register(s)

	address := "0.0.0.0:5060"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Backend Error %v", err)
	}

	fmt.Println("txnService is listening on", address)

	db.ConnectDB()
	defer db.CloseDB()

	// Wait for a single response
	subscriber := redisClient.Subscribe(GET_USER_BALANCE_TOPIC)

	go func() {
		for {
			select {
			case msg := <-subscriber.Channel():
				fmt.Println("TXN - Received message from " + msg.Channel + " channel.")
				fmt.Println(" TXN - payload " + msg.Payload)

				txnRepo := repo.New(db.DB)
				svc := service.New(txnRepo)
				balance, err := svc.GetUserBalance(msg.Payload)
				if err != nil {
					log.Fatalf("Backend Error %v", err)
				}

				if err := redisClient.Publish(USER_BALANCE_TOPIC, balance).Err(); err != nil {
					panic(err)
				}
			}
		}
	}()

	s.Serve(lis)
}
