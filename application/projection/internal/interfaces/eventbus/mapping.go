package eventbus

import providermodel "github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
import "github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"

func projectionUser(user providermodel.User) model.User {
	return model.User{
		Id:    model.UserId(user.Id),
		Email: model.Email(user.Email),
	}
}
