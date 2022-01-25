package repository

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"io"
)

type UserRepository interface {
	io.Closer
	AddUser(ctx context.Context, user model.User) error
	DeleteUserById(ctx context.Context, userId model.UserId) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	FindUserByEmail(ctx context.Context, email model.Email) (model.User, error)
}
