package http

import (
	"github.com/frederic-gendebien/poc-pact/server/internal/usecase"
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
