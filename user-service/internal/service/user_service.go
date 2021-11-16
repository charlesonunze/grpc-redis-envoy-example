package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/repo"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
)

// UserService - interface for the user service
type UserService interface {
	LoginUser(ctx context.Context, name string) (string, error)
	GetUserBalance(ctx context.Context, tkn string) (int32, error)
}

type userService struct {
	repo repo.UserRepo
}

// New - returns an instance of the UserService
func New(repo repo.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

var (
	mySigningKey = []byte(os.Getenv("SECRET_KEY"))
	redisClient  = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PWD"),
	})
	GET_USER_BALANCE_TOPIC = "get_user_balance"
	USER_BALANCE_TOPIC     = "user_balance"
)

// LoginUser - logs in the user and returns a jwt
func (s *userService) LoginUser(ctx context.Context, name string) (string, error) {
	user, err := s.repo.GetUserByName(name)
	if err != nil {
		return "", err
	}

	if user.ID == "" {
		return "", nil
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserBalance - returns the user balances
func (s *userService) GetUserBalance(ctx context.Context, tkn string) (int32, error) {
	var balance int32

	// verify token
	userID, err := verifyToken(tkn)
	if err != nil {
		return balance, err
	}

	// Wait for a single response
	subscriber := redisClient.Subscribe(USER_BALANCE_TOPIC)

	// publish user id in topic
	if err := redisClient.Publish(GET_USER_BALANCE_TOPIC, userID).Err(); err != nil {
		panic(err)
	}

	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		fmt.Println("USR - Received message from " + msg.Channel + " channel.")
		fmt.Println("USR - payload " + msg.Payload)

		i, err := strconv.ParseInt(msg.Payload, 10, 32)
		if err != nil {
			panic(err)
		}

		balance = int32(i)
		break
	}

	return balance, nil
}

func verifyToken(tkn string) (string, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(("Invalid Signing Method"))
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			return nil, fmt.Errorf(("Expired token"))
		}

		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["id"]), nil
	}

	return "", fmt.Errorf(("Token verification failed"))
}

func generateJWT(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("Something Went Wrong: %s", err.Error())
	}

	return tokenString, nil
}
