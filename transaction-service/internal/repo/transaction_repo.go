package repo

import (
	"github.com/charlesonunze/grpc-redis-envoy-example/transaction-service/internal/model"
	"gorm.io/gorm"
)

// TxnRepo - transaction repository
type TxnRepo interface {
	CreditUserAccount(userID string, amount int32) error
	DebitUserAccount(userID string, amount int32) error
	GetUserBalance(userID string) (int32, error)
}

type txnRepo struct {
	db *gorm.DB
}

// New - returns a new transaction repo
func New(db *gorm.DB) TxnRepo {
	return &txnRepo{
		db: db,
	}
}

func (r *txnRepo) GetUserBalance(userID string) (int32, error) {
	var account model.Account
	if err := r.db.Find(&account, "user_id = ?", userID).Error; err != nil {
		return 0, err
	}

	return int32(account.Balance), nil
}

func (r *txnRepo) CreditUserAccount(userID string, amount int32) error {
	if err := r.db.Model(&model.Account{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return err
	}

	return nil
}

func (r *txnRepo) DebitUserAccount(userID string, amount int32) error {
	if err := r.db.Model(&model.Account{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		return err
	}

	return nil
}
