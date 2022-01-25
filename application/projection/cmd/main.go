package main

import (
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/infrastructure/persistence"
	handlers "github.com/frederic-gendebien/pact-poc/application/projection/internal/interfaces/eventbus"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/interfaces/http"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	"log"
)

var (
	configuration config.Configuration
	repo          repository.UserRepository
	eventBus      eventbus.EventBus
	useCase       usecase.UserProjectionUseCase
	server        *http.Server
)

func init() {
	configuration = config.NewConfiguration()
	repo = persistence.NewUserRepository(configuration)
	eventBus = eventbus.NewEventBus(configuration)
	useCase = usecase.NewUserProjectionUseCase(repo)
	server = http.NewServer(useCase)
}

func main() {
	defer teardown()

	go func() {
		log.Println("start consuming events")
		if err := eventBus.Listen(
			handlers.NewUserRegisteredHandler(useCase),
			handlers.UserDeletedHandler(useCase),
		); err != nil {
			log.Fatalln("could not listen for events: ", err)
		}
	}()

	log.Println("starting server...")
	log.Fatalln(server.Start())
}

func teardown() {
	log.Println("tearing down server resources")
	_ = repo.Close()
	_ = eventBus.Close()
	_ = configuration.Close()
}
