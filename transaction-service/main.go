package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/db"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/repo"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/rpc"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/service"
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	GET_USER_BALANCE_TOPIC = "get_user_balance"
	USER_BALANCE_TOPIC     = "user_balance"
)

func main() {
	// Connect to NATS
	opts := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   5,
		ReconnectWait:  5 * time.Second,
		Timeout:        3 * time.Second,
		Url:            os.Getenv("NATS_URL"),
	}

	nc, err := opts.Connect()
	if err != nil {
		fmt.Printf("err => %v", err)
		log.Fatal(err)
	}
	defer nc.Close()
	fmt.Println("nats connected")

	db.ConnectDB()
	defer db.CloseDB()

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

	// Channel Subscriber
	userBalanceChan := make(chan *nats.Msg)
	_, err = nc.ChanSubscribe(GET_USER_BALANCE_TOPIC, userBalanceChan)
	if err != nil {
		fmt.Printf("err%v", err)
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case msg := <-userBalanceChan:
				data := string(msg.Data)
				fmt.Println("TXN - Received message from " + msg.Subject + " channel.")
				fmt.Println(" TXN - payload " + data)

				txnRepo := repo.New(db.DB)
				svc := service.New(txnRepo)
				balance, err := svc.GetUserBalance(data)
				if err != nil {
					log.Fatalf("Backend Error %v", err)
				}

				nc.Publish(USER_BALANCE_TOPIC, []byte(strconv.Itoa(int(balance))))
			}
		}
	}()

	s.Serve(lis)
}
