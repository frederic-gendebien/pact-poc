package usecase

import (
	"context"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/domain/repository"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
)

type UserUseCase interface {
	RegisterNewUser(ctx context.Context, newUser model.User) error
	CorrectUserDetails(ctx context.Context, userId model.UserId, newDetails model.UserDetails) error
	DeleteUser(ctx context.Context, userId model.UserId) error
	ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error)
	FindUserById(ctx context.Context, userId model.UserId) (model.User, error)
}

func NewUserUseCase(repository repository.UserRepository, eventBus eventbus.EventBus) *DefaultUserUseCase {
	return &DefaultUserUseCase{
		repository: repository,
		eventBus:   eventBus,
	}
}

type DefaultUserUseCase struct {
	repository repository.UserRepository
	eventBus   eventbus.EventBus
}

func (d *DefaultUserUseCase) RegisterNewUser(ctx context.Context, newUser model.User) error {
	if err := d.repository.AddUser(ctx, newUser); err != nil {
		return err
	}

	return d.eventBus.Publish(ctx, events.NewUserRegistered{
		User: newUser,
	})
}

func (d *DefaultUserUseCase) CorrectUserDetails(ctx context.Context, userId model.UserId, newDetails model.UserDetails) error {
	correctDetails := func(user model.User) model.User {
		return user.CorrectDetails(newDetails)
	}

	if err := d.repository.UpdateUser(ctx, userId, correctDetails); err != nil {
		return err
	}

	return d.eventBus.Publish(ctx, events.UserDetailsCorrected{
		UserId:         userId,
		NewUserDetails: newDetails,
	})
}

func (d *DefaultUserUseCase) DeleteUser(ctx context.Context, userId model.UserId) error {
	if err := d.repository.DeleteUser(ctx, userId); err != nil {
		return err
	}

	return d.eventBus.Publish(ctx, events.UserDeleted{
		UserId: userId,
	})
}

func (d *DefaultUserUseCase) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	return d.repository.ListAllUsers(ctx, next)
}

func (d *DefaultUserUseCase) FindUserById(ctx context.Context, userId model.UserId) (model.User, error) {
	return d.repository.GetUser(ctx, userId)
}
