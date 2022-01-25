package http

import (
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(useCase usecase.UserProjectionUseCase) *Server {
	engine := gin.Default()
	addUserHandlers(engine, useCase)

	return &Server{
		engine: engine,
	}
}

func (s *Server) Start() error {
	return s.engine.Run()
}
