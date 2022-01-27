package events

import (
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
)

type NewUserRegistered struct {
	User model.User `json:"user"`
}

func (n NewUserRegistered) GetDomain() string {
	return Domain
}

func (n NewUserRegistered) GetName() string {
	return "NewUserRegistered"
}

func (n NewUserRegistered) GetType() interface{} {
	return &NewUserRegistered{}
}

func (n NewUserRegistered) GetDefinition() domain.EventDefinition {
	return n
}

func (n NewUserRegistered) GetEntityId() string {
	return string(n.User.Id)
}

func (n NewUserRegistered) GetPayload() interface{} {
	return n
}
