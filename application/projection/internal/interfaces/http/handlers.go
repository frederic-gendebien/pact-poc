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
	users.GET("", findUsers(useCase))
}

func findUsers(useCase usecase.UserProjectionUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if text := ctx.Query("text"); text != "" {
			users, err := useCase.FindUsersByText(ctx, text)
			okOrFail(ctx, err, func() interface{} {
				return users
			})
		}

		fail(ctx, model.NewBadRequest("missing mandatory 'text' query param"))
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
