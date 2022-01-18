package http

import (
	"errors"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"github.com/gin-gonic/gin"
	gohttp "net/http"
)

func NewErrorResponseFrom(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func createdOrFail(ctx *gin.Context, err error) {
	failOrStatus(ctx, err, gohttp.StatusCreated)
}

func acceptedOrFail(ctx *gin.Context, err error) {
	failOrStatus(ctx, err, gohttp.StatusAccepted)
}

func noContentOrFail(ctx *gin.Context, err error) {
	failOrStatus(ctx, err, gohttp.StatusNoContent)
}

func okOrFail(ctx *gin.Context, err error, body func() interface{}) {
	if err != nil {
		fail(ctx, err)
		return
	}

	ctx.JSON(gohttp.StatusOK, body())
}

func failOrStatus(ctx *gin.Context, err error, status int) {
	if err != nil {
		fail(ctx, err)
		return
	}

	ctx.Status(status)
}

func fail(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, model.BadRequestError{}):
		ctx.AbortWithStatusJSON(gohttp.StatusBadRequest, NewErrorResponseFrom(err))
		break
	case errors.Is(err, model.NotFoundError{}):
		ctx.AbortWithStatusJSON(gohttp.StatusNotFound, NewErrorResponseFrom(err))
		break
	case errors.Is(err, model.UnknownError{}):
		ctx.AbortWithStatusJSON(gohttp.StatusInternalServerError, NewErrorResponseFrom(err))
		break
	default:
		ctx.AbortWithStatusJSON(gohttp.StatusInternalServerError, NewErrorResponseFrom(err))
		break
	}
}
