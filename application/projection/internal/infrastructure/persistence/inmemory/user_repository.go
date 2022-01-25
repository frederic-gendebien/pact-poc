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
		lock:   &sync.RWMutex{},
		emails: make(map[model.UserId]model.Email),
		users:  make(map[model.Email]model.User),
	}
}

type UserRepository struct {
	lock   *sync.RWMutex
	emails map[model.UserId]model.Email
	users  map[model.Email]model.User
}

func (u *UserRepository) Close() error {
	log.Println("closing inmemory user repository")

	return nil
}

func (u *UserRepository) AddUser(ctx context.Context, user model.User) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.emails[user.Id] = user.Email
	u.users[user.Email] = user

	return nil
}

func (u *UserRepository) DeleteUserById(ctx context.Context, userId model.UserId) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	if email, present := u.emails[userId]; present {
		delete(u.users, email)
		delete(u.emails, userId)

		return nil
	}

	return model.NewNotFoundError(fmt.Sprintf("user with id: %s was not found", userId))
}

func (u *UserRepository) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	users := make(chan model.User)
	go func() {
		u.lock.RLock()
		defer u.lock.RUnlock()
		defer close(users)

		for _, user := range u.users {
			users <- user
			select {
			case needNext := <-next:
				if !needNext {
					break
				}
			}
		}
	}()

	return users, nil
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email model.Email) (model.User, error) {
	u.lock.RLock()
	defer u.lock.RUnlock()

	if user, present := u.users[email]; present {
		return user, nil
	}

	return model.User{}, model.NewNotFoundError(fmt.Sprintf("user with email: %s was not found", email))
}
