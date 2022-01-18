package domain

import (
	"context"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"io"
)

type UserRepository interface {
	io.Closer
	AddUser(ctx context.Context, newUser model.User) error
	DeleteUser(ctx context.Context, userId string) error
	GetUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	GetUser(ctx context.Context, id string) (model.User, error)
}
