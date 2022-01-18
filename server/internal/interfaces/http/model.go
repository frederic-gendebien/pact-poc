package http

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/pkg/domain/model"
)

func NewUserFrom(user model.User) User {
	return User{
		Id:    user.GetId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) GetId() string {
	return u.Id
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetEmail() string {
	return u.Email
}

func (u User) IsInvalid() bool {
	return u.Id == "" ||
		u.Name == "" ||
		u.Email == ""
}

func (u User) Invalid() error {
	if u.IsInvalid() {
		return model.NewBadRequest("user content is invalid")
	}

	return nil
}

func (u User) InvalidAfter(err error) error {
	if err != nil {
		return err
	}

	return u.Invalid()
}
