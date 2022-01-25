package http

import (
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	MaxLimit int = 20
)

func addUserHandlers(engine *gin.Engine, useCase usecase.UserProjectionUseCase) {
	users := engine.Group("/users")
	users.GET("", listOrFindUsers(useCase))
}

func listOrFindUsers(useCase usecase.UserProjectionUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if email := ctx.Query("email"); email != "" {
			user, err := useCase.FindUserByEmail(ctx, model.Email(email))
			okOrFail(ctx, err, func() interface{} {
				return user
			})
		} else {
			limit := minOrDefault(ctx.Query("limit"), MaxLimit)
			next := make(chan bool)
			defer close(next)

			users, err := useCase.ListAllUsers(ctx, next)

			okOrFail(ctx, err, accumulateUsers(users, limit, next))
		}
	}
}

func minOrDefault(value string, maxValue int) int {
	if value == "" {
		return maxValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return maxValue
	}

	if intValue > MaxLimit {
		return MaxLimit
	}

	return intValue
}

func accumulateUsers(users <-chan model.User, limit int, next chan<- bool) func() interface{} {
	return func() interface{} {
		counter := 0
		results := make([]model.User, 0, MaxLimit)

		for user := range users {
			if counter >= limit {
				next <- false
				break
			} else {
				results = append(results, user)
				counter++
				next <- true
			}
		}

		return results
	}
}
