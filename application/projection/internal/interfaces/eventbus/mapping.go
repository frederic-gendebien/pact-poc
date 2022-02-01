package eventbus

import (
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	providermodel "github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
)
import "github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"

func projectionUser(user providermodel.User) model.User {
	return model.User{
		Id:    model.UserId(user.Id),
		Name:  user.Details.Name,
		Email: model.Email(user.Email),
	}
}

func partialUserFrom(detailsCorrected *events.UserDetailsCorrected) model.User {
	return model.User{
		Id:    model.UserId(detailsCorrected.UserId),
		Name:  detailsCorrected.NewUserDetails.Name,
		Email: "",
	}
}
