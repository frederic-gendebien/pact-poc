package events

import (
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
)

type UserDeleted struct {
	UserId model.UserId `json:"user_id"`
}

func (u UserDeleted) GetDomain() string {
	return Domain
}

func (u UserDeleted) GetName() string {
	return "UserDeleted"
}

func (u UserDeleted) GetDefinition() domain.EventDefinition {
	return u
}

func (u UserDeleted) GetEntityId() string {
	return string(u.UserId)
}

func (u UserDeleted) GetPayload() interface{} {
	return u
}
