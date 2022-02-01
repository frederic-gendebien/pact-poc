package eventbus

import (
	"context"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
	inmemorypers "github.com/frederic-gendebien/pact-poc/application/projection/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	providermodel "github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
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
			ExpectsToReceive("a user1 registered event").
			WithContent(&events.NewUserRegistered{User: user}).
			AsType(&events.NewUserRegistered{})

		if err := pact.VerifyMessageConsumer(t, message, sendSearchAndExpect(newUserRegistered(), string(user.Email), user)); err != nil {
			t.Fatalf("Error on verify: %v", err)
		}
	})
}

func TestProjectionPact_UserDetailsCorrected(t *testing.T) {
	t.Run("User Details Corrected", func(t *testing.T) {
		user := testUser(1)
		message := pact.AddMessage()
		message.Given("user1 details have been corrected").
			ExpectsToReceive("a user1 details corrected event").
			WithContent(&events.UserDetailsCorrected{
				UserId:         user.Id,
				NewUserDetails: newTestUserDetails(1),
			}).
			AsType(&events.UserDetailsCorrected{})

		if err := pact.VerifyMessageConsumer(t, message, sendSearchAndExpect(userDetailsCorrected(), string(user.Email), user)); err != nil {
			t.Fatalf("Error on verify: %v", err)
		}
	})
}

func sendSearchAndExpect(eventFrom func(dsl.Message) domain.Event, text string, registeredUser providermodel.User) dsl.MessageConsumer {
	return func(message dsl.Message) error {
		if err := eventBus.Publish(context.Background(), eventFrom(message)); err != nil {
			return err
		}

		persistedUsers, err := useCase.FindUsersByText(context.Background(), text)
		if err != nil {
			return err
		}

		if len(persistedUsers) == 0 {
			return fmt.Errorf("no users found")
		}

		if string(persistedUsers[0].Id) != string(registeredUser.Id) ||
			string(persistedUsers[0].Email) != string(registeredUser.Email) {
			return fmt.Errorf("expected user was: %v, but received: %v", registeredUser, persistedUsers)
		}

		return nil
	}
}

func newUserRegistered() func(dsl.Message) domain.Event {
	return func(message dsl.Message) domain.Event {
		return message.Content.(*events.NewUserRegistered)
	}
}

func userDetailsCorrected() func(dsl.Message) domain.Event {
	return func(message dsl.Message) domain.Event {
		return message.Content.(*events.UserDetailsCorrected)
	}
}

func setup() {
	log.Println("setup pact environment")
	pact = dsl.Pact{
		Consumer:                 "user-projection",
		Provider:                 "user-server-usecase",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
		LogLevel:                 "INFO",
		DisableToolValidityCheck: true,
	}
	pact.Setup(false)
	repo = inmemorypers.NewUserRepository()
	useCase = usecase.NewUserProjectionUseCase(repo)
	eventBus = inmemoryevb.NewEventBus()
	if err := eventBus.Listen(context.Background(),
		ListenerName,
		NewUserRegisteredHandler(useCase),
		UserDetailsCorrectedHandler(useCase),
		UserDeletedHandler(useCase),
	); err != nil {
		log.Fatalf("could not listen for events: %v", err)
	}

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
		Id: providermodel.UserId(fmt.Sprintf("user%d", number)),
		Details: providermodel.UserDetails{
			Name: fmt.Sprintf("name%d", number),
		},
		Email: providermodel.Email(fmt.Sprintf("email%d", number)),
	}
}

func newTestUserDetails(number int) providermodel.UserDetails {
	return providermodel.UserDetails{
		Name: fmt.Sprintf("new_name%d", number),
	}
}
