package events

import (
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
)

type UserDetailsCorrected struct {
	UserId         model.UserId      `json:"user_id"`
	NewUserDetails model.UserDetails `json:"new_user_details"`
}

func (n UserDetailsCorrected) GetDomain() string {
	return Domain
}

func (n UserDetailsCorrected) GetName() string {
	return "UserDetailsCorrected"
}

func (n UserDetailsCorrected) GetType() interface{} {
	return &UserDetailsCorrected{}
}

func (n UserDetailsCorrected) GetDefinition() domain.EventDefinition {
	return n
}

func (n UserDetailsCorrected) GetEntityId() string {
	return string(n.UserId)
}

func (n UserDetailsCorrected) GetPayload() interface{} {
	return n
}
