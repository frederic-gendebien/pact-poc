package main

import (
	"github.com/frederic-gendebien/pact-poc/application/server/internal/domain"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/interfaces/http"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/usecase"
	"log"
)

var (
	repository domain.UserRepository
	useCase    usecase.UserUseCase
	server     *http.Server
)

func init() {
	repository = inmemory.NewUserRepository()
	useCase = usecase.NewUserUseCase(repository)
	server = http.NewServer(useCase)
}

func main() {
	defer teardown()

	log.Println("starting server...")
	log.Fatalln(server.Start())
}

func teardown() {
	log.Println("tearing down server resources")
	_ = repository.Close()
}
