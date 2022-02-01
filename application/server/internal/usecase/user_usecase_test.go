package usecase

import (
	"context"
	"fmt"
	inmemorypers "github.com/frederic-gendebien/pact-poc/application/server/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/config/environment"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	inmemoryevb "github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"testing"
)

const (
	pactBrokerUrlPropertyName   = "PACT_BROKER_URL"
	pactBrokerTokenPropertyName = "PACT_BROKER_TOKEN"
)

var (
	configuration   config.Configuration
	pactBrokerUrl   string
	pactBrokerToken string
	repo            *inmemorypers.UserRepository
	eventBus        *inmemoryevb.EventBus
	eventSniffer    *eventbus.EventSniffer
	useCase         UserUseCase
)

func init() {
	configuration = environment.NewConfiguration()
	pactBrokerUrl = configuration.GetStringOrCrash(pactBrokerUrlPropertyName)
	pactBrokerToken = configuration.GetStringOrCrash(pactBrokerTokenPropertyName)

	repo = inmemorypers.NewUserRepository()
	eventBus = inmemoryevb.NewEventBus()
	eventSniffer = eventbus.NewEventSniffer(eventBus)
	useCase = NewUserUseCase(repo, eventBus)
}

func TestServerMessagePact(t *testing.T) {
	eventSniffer.Clear()
	if err := eventSniffer.Listen(
		events.NewUserRegistered{},
		events.UserDetailsCorrected{},
	); err != nil {
		t.Fatal(err)
	}

	pact := dsl.Pact{
		Provider:                 "user-server-usecase",
		LogDir:                   "../../../tests/pact/logs",
		PactDir:                  "../../../tests/pact/pacts",
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}

	if _, err := pact.VerifyMessageProvider(t, dsl.VerifyMessageRequest{
		Tags:                       []string{"main"},
		BrokerURL:                  pactBrokerUrl,
		BrokerToken:                pactBrokerToken,
		PactURLs:                   nil,
		ConsumerVersionSelectors:   nil,
		PublishVerificationResults: true,
		ProviderVersion:            "0.0.1",
		ProviderTags:               []string{"main"},
		MessageHandlers:            messageHandlers(),
		StateHandlers:              messageStateHandlers(),
		PactLogDir:                 "../../../../tests/pact/logs",
	}); err != nil {
		t.Fatalf("server message verifaction failed: %v", err)
	}
}

func messageHandlers() dsl.MessageHandlers {
	return dsl.MessageHandlers{
		"a user1 registered event": func(message dsl.Message) (interface{}, error) {
			return eventSniffer.GetAndClearEvents()[0], nil
		},
		"a user1 details corrected event": func(message dsl.Message) (interface{}, error) {
			return eventSniffer.GetAndClearEvents()[0], nil
		},
	}
}

func messageStateHandlers() dsl.StateHandlers {
	return dsl.StateHandlers{
		"user1 has been registered": func(state dsl.State) error {
			return useCase.RegisterNewUser(context.Background(), testUser(1))
		},
		"user1 details have been corrected": func(state dsl.State) error {
			return useCase.CorrectUserDetails(context.Background(), testUser(1).Id, newTestUserDetails(1))
		},
	}
}

func emptyRepository() types.StateHandler {
	return func() error {
		return repo.Clear(context.Background())
	}
}

func repositoryWith(users ...model.User) func() error {
	return func() error {
		ctx := context.Background()
		if err := repo.Clear(ctx); err != nil {
			return err
		}

		for _, user := range users {
			if err := repo.AddUser(ctx, user); err != nil {
				return err
			}
		}

		return nil
	}
}

func testUser(number int) model.User {
	return model.User{
		Id: model.UserId(fmt.Sprintf("user%d", number)),
		Details: model.UserDetails{
			Name: fmt.Sprintf("name%d", number),
		},
		Email: model.Email(fmt.Sprintf("email%d", number)),
	}
}

func newTestUserDetails(number int) model.UserDetails {
	return model.UserDetails{
		Name: fmt.Sprintf("new_name%d", number),
	}
}
