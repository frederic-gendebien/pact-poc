package repository

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"io"
)

type UserRepository interface {
	io.Closer
	AddUser(ctx context.Context, newUser model.User) error
	DeleteUser(ctx context.Context, userId model.UserId) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	GetUser(ctx context.Context, userId model.UserId) (model.User, error)
}
