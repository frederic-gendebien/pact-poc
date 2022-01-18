package http

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(useCase usecase.UserUseCase) *Server {
	engine := gin.Default()
	addUserHandlers(engine, useCase)

	return &Server{
		engine: engine,
	}
}

func (s *Server) Start() error {
	return s.engine.Run()
}
