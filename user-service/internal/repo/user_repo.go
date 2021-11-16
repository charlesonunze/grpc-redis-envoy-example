package repo

import (
	"github.com/charlesonunze/grpc-redis-envoy-example/user-service/internal/model"
	"gorm.io/gorm"
)

// UserRepo - user repository
type UserRepo interface {
	GetUserByName(name string) (model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

// New - returns a new user repo
func New(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetUserByName(name string) (model.User, error) {
	var user model.User

	if err := r.db.Find(&user, "name = ?", name).Error; err != nil {
		return user, err
	}

	return user, nil
}
