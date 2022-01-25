package usecase

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/domain"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
)

type UserUseCase interface {
	RegisterNewUser(ctx context.Context, newUser model.User) error
	DeleteUser(ctx context.Context, userId string) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	FindUserById(ctx context.Context, userId string) (model.User, error)
}

func NewUserUseCase(repository domain.UserRepository, eventBus eventbus.EventBus) *DefaultUserCase {
	return &DefaultUserCase{
		repository: repository,
		eventBus:   eventBus,
	}
}

type DefaultUserCase struct {
	repository domain.UserRepository
	eventBus   eventbus.EventBus
}

func (d *DefaultUserCase) RegisterNewUser(ctx context.Context, newUser model.User) error {
	if err := d.repository.AddUser(ctx, newUser); err != nil {
		return err
	}

	return d.eventBus.Publish(events.NewUserRegistered{
		User: newUser,
	})
}

func (d *DefaultUserCase) DeleteUser(ctx context.Context, userId string) error {
	if err := d.repository.DeleteUser(ctx, userId); err != nil {
		return err
	}

	return d.eventBus.Publish(events.UserDeleted{
		UserId: userId,
	})
}

func (d *DefaultUserCase) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	return d.repository.GetUsers(ctx, next)
}

func (d *DefaultUserCase) FindUserById(ctx context.Context, userId string) (model.User, error) {
	return d.repository.GetUser(ctx, userId)
}
