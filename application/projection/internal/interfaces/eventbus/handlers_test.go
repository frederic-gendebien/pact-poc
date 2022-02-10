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
	"github.com/pact-foundation/pact-go/v2/message"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	pact     *message.MessagePactV3
	repo     repository.UserRepository
	useCase  usecase.UserProjectionUseCase
	eventBus eventbus.EventBus
)

func init() {
	log.Println("setup pact environment")
	var err error
	pact, err = message.NewMessagePactV3(message.MessageConfig{
		Consumer: "user-projection",
		Provider: "user-server-usecase",
		PactDir:  "../../../../tests/pact/pacts",
	})
	if err != nil {
		log.Fatalf("could not listen for events: %v", err)
	}

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

func TestProjectionPact_NewUserRegistered(t *testing.T) {
	t.Run("New User Registered", func(t *testing.T) {
		user := testUser(1)
		err := pact.AddMessage().
			Given(models.ProviderStateV3{
				Name:       "user1 has been registered",
				Parameters: userParameters(user),
			}).
			ExpectsToReceive("a user1 registered event").
			WithJSONContent(&events.NewUserRegistered{User: user}).
			AsType(&events.NewUserRegistered{}).
			ConsumedBy(sendSearchAndExpect(newUserRegistered(), string(user.Email), user)).
			Verify(t)

		assert.NoError(t, err)
	})
}

func TestProjectionPact_UserDetailsCorrected(t *testing.T) {
	t.Run("User Details Corrected", func(t *testing.T) {
		user := testUser(1)
		err := pact.AddMessage().
			Given(models.ProviderStateV3{
				Name:       "user1 details have been corrected",
				Parameters: userParameters(user),
			}).
			ExpectsToReceive("a user1 details corrected event").
			WithJSONContent(&events.UserDetailsCorrected{
				UserId:         user.Id,
				NewUserDetails: newTestUserDetails(1),
			}).
			AsType(&events.UserDetailsCorrected{}).
			ConsumedBy(sendSearchAndExpect(userDetailsCorrected(), string(user.Email), user)).
			Verify(t)

		assert.NoError(t, err)
	})
}

func userParameters(user providermodel.User) map[string]interface{} {
	return map[string]interface{}{
		"Id": user.Id,
	}
}

func sendSearchAndExpect(eventFrom func(message.AsynchronousMessage) domain.Event, text string, registeredUser providermodel.User) message.MessageConsumer {
	return func(asynchronousMessage message.AsynchronousMessage) error {
		if err := eventBus.Publish(context.Background(), eventFrom(asynchronousMessage)); err != nil {
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

func newUserRegistered() func(message.AsynchronousMessage) domain.Event {
	return func(event message.AsynchronousMessage) domain.Event {
		return event.Content.(*events.NewUserRegistered)
	}
}

func userDetailsCorrected() func(message.AsynchronousMessage) domain.Event {
	return func(event message.AsynchronousMessage) domain.Event {
		return event.Content.(*events.UserDetailsCorrected)
	}
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
