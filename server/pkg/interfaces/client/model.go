package client

import "github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"

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
