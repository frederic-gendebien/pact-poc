package usecase

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/domain"
	"bitbucket.org/fredericgendebien/pact-poc/server/pkg/domain/model"
	"context"
)

type UserUseCase interface {
	RegisterNewUser(ctx context.Context, newUser model.User) error
	DeleteUser(ctx context.Context, userId string) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	FindUserById(ctx context.Context, userId string) (model.User, error)
}

func NewUserUseCase(repository domain.UserRepository) *DefaultUserCase {
	return &DefaultUserCase{
		repository: repository,
	}
}

type DefaultUserCase struct {
	repository domain.UserRepository
}

func (d *DefaultUserCase) RegisterNewUser(ctx context.Context, newUser model.User) error {
	return d.repository.AddUser(ctx, newUser)
}

func (d *DefaultUserCase) DeleteUser(ctx context.Context, userId string) error {
	return d.repository.DeleteUser(ctx, userId)
}

func (d *DefaultUserCase) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	return d.repository.GetUsers(ctx, next)
}

func (d *DefaultUserCase) FindUserById(ctx context.Context, userId string) (model.User, error) {
	return d.repository.GetUser(ctx, userId)
}
