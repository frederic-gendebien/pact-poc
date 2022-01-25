package eventbus

import (
	"context"
	"fmt"
	model "github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/model"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
	inmemorypers "github.com/frederic-gendebien/pact-poc/application/projection/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	providermodel "github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	inmemoryevb "github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"github.com/pact-foundation/pact-go/dsl"
	"log"
	"os"
	"testing"
)

var (
	pact     dsl.Pact
	repo     repository.UserRepository
	useCase  usecase.UserProjectionUseCase
	eventBus eventbus.EventBus
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		setup()
		defer tearDown()

		exitCode := m.Run()

		if err := pact.WritePact(); err != nil {
			log.Fatalln("could not write pact: ", err)
		}

		return exitCode
	}())
}

func TestProjectionPact_NewUserRegistered(t *testing.T) {
	t.Run("New User Registered", func(t *testing.T) {
		user := testUser(1)
		message := pact.AddMessage()
		message.Given("user1 has been registered").
			ExpectsToReceive("a new user registered event").
			WithContent(map[string]interface{}{
				"id":    user.Id,
				"name":  user.Name,
				"email": user.Email,
			}).
			AsType(&events.NewUserRegistered{User: user})

		if err := pact.VerifyMessageConsumer(t, message, expectUser(user)); err != nil {
			t.Fatalf("Error on verify: %v", err)
		}
	})
}

func expectUser(registeredUser providermodel.User) dsl.MessageConsumer {
	return func(message dsl.Message) error {
		event := message.Content.(*events.NewUserRegistered)
		if err := eventBus.Listen(
			NewUserRegisteredHandler(useCase),
			UserDeletedHandler(useCase),
		); err != nil {
			return err
		}

		if err := eventBus.Publish(event); err != nil {
			return err
		}

		persistedUser, err := useCase.FindUserByEmail(context.Background(), model.Email(event.User.Email))
		if err != nil {
			return err
		}

		if string(persistedUser.Id) != string(registeredUser.Id) ||
			string(persistedUser.Email) != string(registeredUser.Email) {
			return fmt.Errorf("expected user was: %v, but received: %v", registeredUser, persistedUser)
		}

		return nil
	}
}

func setup() {
	log.Println("setup pact environment")
	pact = dsl.Pact{
		Consumer:                 "user-projection",
		Provider:                 "user-server",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
		LogLevel:                 "INFO",
		DisableToolValidityCheck: true,
	}
	pact.Setup(false)
	repo = inmemorypers.NewUserRepository()
	useCase = usecase.NewUserProjectionUseCase(repo)
	eventBus = inmemoryevb.NewEventBus()

	log.Println("pact environment setup done")
}

func tearDown() {
	log.Println("tearing down pact environment")
	pact.Teardown()
	_ = repo.Close()
	_ = eventBus.Close()
	log.Println("pact environment tear down done")
}

func testUser(number int) providermodel.User {
	return providermodel.User{
		Id:    providermodel.UserId(fmt.Sprintf("user%d", number)),
		Name:  fmt.Sprintf("name%d", number),
		Email: providermodel.Email(fmt.Sprintf("email%d", number)),
	}
}
