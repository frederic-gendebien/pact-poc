package main

import (
	"github.com/frederic-gendebien/pact-poc/application/server/internal/domain"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/infrastructure/persistence"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/interfaces/http"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	"log"
)

var (
	configuration config.Configuration
	repository    domain.UserRepository
	eventBus      eventbus.EventBus
	useCase       usecase.UserUseCase
	server        *http.Server
)

func init() {
	configuration = config.NewConfiguration()
	repository = persistence.NewUserRepository(configuration)
	eventBus = eventbus.NewEventBus(configuration)
	useCase = usecase.NewUserUseCase(repository, eventBus)
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
	_ = eventBus.Close()
	_ = configuration.Close()
}
