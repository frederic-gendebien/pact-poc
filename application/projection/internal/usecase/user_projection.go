package usecase

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
)

type UserProjectionUseCase interface {
	IndexUser(ctx context.Context, user model.User) error
	DeleteUserById(ctx context.Context, userId model.UserId) error
	FindUsersByText(ctx context.Context, text string) ([]model.User, error)
}

func NewUserProjectionUseCase(repository repository.UserRepository) *DefaultUserProjectionUseCase {
	return &DefaultUserProjectionUseCase{
		repository: repository,
	}
}

type DefaultUserProjectionUseCase struct {
	repository repository.UserRepository
}

func (d *DefaultUserProjectionUseCase) IndexUser(ctx context.Context, user model.User) error {
	return d.repository.IndexUser(ctx, user)
}

func (d *DefaultUserProjectionUseCase) DeleteUserById(ctx context.Context, userId model.UserId) error {
	return d.repository.DeleteUserById(ctx, userId)
}

func (d *DefaultUserProjectionUseCase) FindUsersByText(ctx context.Context, text string) ([]model.User, error) {
	return d.repository.FindUsersByText(ctx, text)
}
