package inmemory

import (
	"context"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"log"
	"sync"
)

func NewUserRepository() *UserRepository {
	return &UserRepository{
		lock:     &sync.RWMutex{},
		patterns: make(map[string][]model.UserId),
		users:    make(map[model.UserId]model.User),
	}
}

type UserRepository struct {
	lock     *sync.RWMutex
	patterns map[string][]model.UserId
	users    map[model.UserId]model.User
}

func (u *UserRepository) Close() error {
	log.Println("closing inmemory user repository")

	return nil
}

func (u *UserRepository) IndexUser(ctx context.Context, user model.User) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	resultingUser, present := u.users[user.Id]
	if present {
		resultingUser = resultingUser.UpdateWith(user)
	} else {
		resultingUser = user
	}

	u.users[resultingUser.Id] = resultingUser
	u.addPattern(ctx, resultingUser.Id, resultingUser.Name)
	u.addPattern(ctx, resultingUser.Id, string(resultingUser.Email))

	return nil
}

func (u *UserRepository) addPattern(ctx context.Context, userId model.UserId, text string) {
	if text != "" {
		u.patterns[text] = append(u.patterns[text], userId)
	}
}

func (u *UserRepository) DeleteUserById(ctx context.Context, userId model.UserId) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	if _, present := u.users[userId]; present {
		delete(u.users, userId)

		return nil
	}

	return model.NewNotFoundError(fmt.Sprintf("user with id: %s was not found", userId))
}

func (u *UserRepository) FindUsersByText(ctx context.Context, text string) ([]model.User, error) {
	u.lock.RLock()
	defer u.lock.RUnlock()

	users := make([]model.User, 0, 2)
	if userIds := u.patterns[text]; userIds != nil {
		for _, userId := range userIds {
			if user, present := u.users[userId]; present {
				users = append(users, user)
			}
		}
	}

	return users, nil
}
