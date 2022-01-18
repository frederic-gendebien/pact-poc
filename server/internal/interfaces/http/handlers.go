package http

import (
	"github.com/frederic-gendebien/poc-pact/server/internal/usecase"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	MaxLimit int = 20
)

func addUserHandlers(engine *gin.Engine, useCase usecase.UserUseCase) {
	users := engine.Group("/users")
	users.PUT("", registerNewUser(useCase))
	users.DELETE(":user_id", deleteUser(useCase))
	users.GET("", getUsers(useCase))
	users.GET(":user_id", getUser(useCase))
}

func registerNewUser(useCase usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newUser := User{}
		if err := newUser.InvalidAfter(ctx.BindJSON(&newUser)); err != nil {
			fail(ctx, model.NewBadRequest(err.Error()))
			return
		}

		createdOrFail(ctx, useCase.RegisterNewUser(ctx, newUser))
	}
}

func deleteUser(useCase usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")
		if userId == "" {
			fail(ctx, model.NewBadRequest("wrong user id"))
			return
		}

		acceptedOrFail(ctx, useCase.DeleteUser(ctx, userId))
	}
}

func getUsers(useCase usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit := minOrDefault(ctx.Query("limit"), MaxLimit)
		next := make(chan bool)
		defer close(next)

		users, err := useCase.ListAllUsers(ctx, next)

		okOrFail(ctx, err, accumulateUsers(users, limit, next))
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
		results := make([]User, 0, MaxLimit)

		for user := range users {
			if counter >= limit {
				next <- false
				break
			} else {
				results = append(results, NewUserFrom(user))
				counter++
				next <- true
			}
		}

		return results
	}
}

func getUser(useCase usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")
		if userId == "" {
			fail(ctx, model.NewBadRequest("wrong user id"))
			return
		}

		user, err := useCase.FindUserById(ctx, userId)

		okOrFail(ctx, err, lazyUserConversion(user))
	}
}

func lazyUserConversion(user model.User) func() interface{} {
	return func() interface{} {
		return NewUserFrom(user)
	}
}
