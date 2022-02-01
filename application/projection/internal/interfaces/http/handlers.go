package http

import (
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/gin-gonic/gin"
)

func addUserHandlers(engine *gin.Engine, useCase usecase.UserProjectionUseCase) {
	users := engine.Group("/users")
	users.GET("", findUsers(useCase))
}

func findUsers(useCase usecase.UserProjectionUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		text := ctx.Query("text")
		if text == "" {
			fail(ctx, model.NewBadRequest("missing mandatory 'text' query param"))
			return
		}

		users, err := useCase.FindUsersByText(ctx, text)
		okOrFail(ctx, err, func() interface{} {
			return users
		})
	}
}
