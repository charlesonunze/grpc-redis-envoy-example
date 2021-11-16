package service

import (
	"context"
	"fmt"
	"os"

	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/repo"
	"github.com/golang-jwt/jwt"
)

// TransactionService - interface for the transaction service
type TransactionService interface {
	CreditUserAccount(ctx context.Context, token string, amount int32) error
	DebitUserAccount(ctx context.Context, token string, amount int32) error
	GetUserBalance(userID string) (int32, error)
}

type transactionService struct {
	repo repo.TxnRepo
}

// New - returns an instance of the TransactionService
func New(repo repo.TxnRepo) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

var (
	mySigningKey = []byte(os.Getenv("SECRET_KEY"))
)

// GetUserBalance - returns the balance of a particular user
func (s *transactionService) GetUserBalance(userID string) (int32, error) {
	balance, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

// CreditUserAccount - adds amount to the user balance
func (s *transactionService) CreditUserAccount(ctx context.Context, token string, amount int32) error {
	// verify token
	userID, err := verifyToken(token)
	if err != nil {
		return err
	}

	err = s.repo.CreditUserAccount(userID, amount)
	if err != nil {
		return err
	}

	return nil
}

// DebitUserAccount - subtracts amount from the user balance
func (s *transactionService) DebitUserAccount(ctx context.Context, token string, amount int32) error {
	// verify token
	userID, err := verifyToken(token)
	if err != nil {
		return err
	}

	err = s.repo.DebitUserAccount(userID, amount)
	if err != nil {
		return err
	}

	return nil
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
