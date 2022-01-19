package main

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/domain"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/infrastructure/persistence/inmemory"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/interfaces/http"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/usecase"
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
