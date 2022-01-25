package usecase

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
)

type UserProjectionUseCase interface {
	AddUser(ctx context.Context, user model.User) error
	DeleteUserById(ctx context.Context, userId model.UserId) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	FindUserByEmail(ctx context.Context, email model.Email) (model.User, error)
}

func NewUserProjectionUseCase(repository repository.UserRepository) *DefaultUserProjectionUseCase {
	return &DefaultUserProjectionUseCase{
		repository: repository,
	}
}

type DefaultUserProjectionUseCase struct {
	repository repository.UserRepository
}

func (d *DefaultUserProjectionUseCase) AddUser(ctx context.Context, user model.User) error {
	return d.repository.AddUser(ctx, user)
}

func (d *DefaultUserProjectionUseCase) DeleteUserById(ctx context.Context, userId model.UserId) error {
	return d.repository.DeleteUserById(ctx, userId)
}

func (d *DefaultUserProjectionUseCase) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	return d.repository.ListAllUsers(ctx, next)
}

func (d *DefaultUserProjectionUseCase) FindUserByEmail(ctx context.Context, email model.Email) (model.User, error) {
	return d.repository.FindUserByEmail(ctx, email)
}
