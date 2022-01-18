package inmemory

import (
	"context"
	"fmt"
	"github.com/frederic-gendebien/poc-pact/server/internal/domain"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"log"
	"sync"
)

func NewUserRepository() domain.UserRepository {
	return &Repository{
		lock:  &sync.RWMutex{},
		users: make(map[string]model.User),
	}
}

type Repository struct {
	lock  *sync.RWMutex
	users map[string]model.User
}

func (r *Repository) Close() error {
	log.Println("closing user repository")

	return nil
}

func (r *Repository) AddUser(ctx context.Context, newUser model.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	userId := newUser.GetId()
	if _, present := r.users[userId]; present {
		return model.NewBadRequest(fmt.Sprintf("user with id: %s already exists", userId))
	}

	r.users[userId] = newUser

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userId string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, present := r.users[userId]; !present {
		return notFound(userId)
	}

	delete(r.users, userId)

	return nil
}

func (r *Repository) GetUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	users := make(chan model.User)
	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()
		defer close(users)

		for _, user := range r.users {
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

func (r *Repository) GetUser(ctx context.Context, userId string) (model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	user, present := r.users[userId]
	if !present {
		return nil, notFound(userId)
	}

	return user, nil
}

func notFound(userId string) model.NotFoundError {
	return model.NewNotFoundError(fmt.Sprintf("user with id: %s was not found", userId))
}
