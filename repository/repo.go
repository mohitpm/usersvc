package repository

import (
	"context"

	"github.com/mohitpm/usersvc/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (repo *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return nil
}
