package inmemory

import (
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"sort"
)

type OrderedMap map[string]model.User

func (om OrderedMap) OrderedValues() []model.User {
	if om == nil {
		return []model.User{}
	}

	var users []model.User
	for _, user := range om {
		users = append(users, user)
	}

	if users == nil {
		return []model.User{}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	return users
}
