package inmemory

import (
	"context"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"log"
	"sync"
)

func NewUserRepository() *UserRepository {
	log.Println("starting inmemory user repository")
	return &UserRepository{
		lock:  &sync.RWMutex{},
		users: make(map[model.UserId]model.User),
	}
}

type UserRepository struct {
	lock  *sync.RWMutex
	users map[model.UserId]model.User
}

func (r *UserRepository) Close() error {
	log.Println("closing inmemory user repository")

	return nil
}

func (r *UserRepository) Clear(ctx context.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users = make(map[model.UserId]model.User)

	return nil
}

func (r *UserRepository) AddUser(ctx context.Context, newUser model.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	userId := newUser.Id
	if _, present := r.users[userId]; present {
		return model.NewBadRequest(fmt.Sprintf("user with id: %s already exists", userId))
	}

	r.users[userId] = newUser

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId model.UserId) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, present := r.users[userId]; !present {
		return notFound(userId)
	}

	delete(r.users, userId)

	return nil
}

func (r *UserRepository) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	users := make(chan model.User)
	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()
		defer close(users)

		for _, user := range OrderedMap(r.users).OrderedValues() {
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

func (r *UserRepository) GetUser(ctx context.Context, userId model.UserId) (model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if user, present := r.users[userId]; present {
		return user, nil
	}

	return model.User{}, notFound(userId)
}

func notFound(userId model.UserId) model.NotFoundError {
	return model.NewNotFoundError(fmt.Sprintf("user with id: %s was not found", userId))
}
