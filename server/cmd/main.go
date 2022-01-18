package main

import (
	"github.com/frederic-gendebien/poc-pact/server/internal/domain"
	"github.com/frederic-gendebien/poc-pact/server/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/poc-pact/server/internal/interfaces/http"
	"github.com/frederic-gendebien/poc-pact/server/internal/usecase"
	"log"
)

var (
	repository domain.UserRepository
	useCase    usecase.UserUseCase
	server     *http.Server
)

func init() {
	repository = inmemory.NewUserRepository()
	useCase = usecase.NewDefaultUserCase(repository)
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
