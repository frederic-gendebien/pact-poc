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
		lock:   &sync.RWMutex{},
		users:  make(map[model.UserId]model.User),
		emails: make(map[model.Email]model.UserId),
	}
}

type UserRepository struct {
	lock   *sync.RWMutex
	users  map[model.UserId]model.User
	emails map[model.Email]model.UserId
}

func (r *UserRepository) Close() error {
	log.Println("closing inmemory user repository")

	return nil
}

func (r *UserRepository) Clear(ctx context.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users = make(map[model.UserId]model.User)
	r.emails = make(map[model.Email]model.UserId)

	return nil
}

func (r *UserRepository) AddUser(ctx context.Context, newUser model.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	email := newUser.Email
	if _, present := r.emails[email]; present {
		return model.NewBadRequest(fmt.Sprintf("user email : %s already exists", email))
	}

	userId := newUser.Id
	if _, present := r.users[userId]; present {
		return model.NewBadRequest(fmt.Sprintf("user with id: %s already exists", userId))
	}

	r.users[userId] = newUser
	r.emails[email] = userId

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId model.UserId) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if user, present := r.users[userId]; present {
		delete(r.users, userId)
		delete(r.emails, user.Email)

		return nil
	}

	return notFound(userId)
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
