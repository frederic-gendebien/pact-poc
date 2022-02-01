package repository

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"io"
)

type UserRepository interface {
	io.Closer
	IndexUser(ctx context.Context, user model.User) error
	DeleteUserById(ctx context.Context, userId model.UserId) error
	FindUsersByText(ctx context.Context, text string) ([]model.User, error)
}
