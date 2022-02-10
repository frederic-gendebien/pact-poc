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
	"github.com/pact-foundation/pact-go/v2/message"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
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
	pact := &message.MessageVerifier{}
	err := pact.Verify(t, message.VerifyMessageRequest{
		VerifyRequest: provider.VerifyRequest{
			BrokerURL:                  pactBrokerUrl,
			BrokerToken:                pactBrokerToken,
			Provider:                   "user-server-usecase",
			PactDirs:                   []string{"../../../tests/pact/pacts"},
			Tags:                       []string{"main"},
			PactURLs:                   nil,
			ConsumerVersionSelectors:   nil,
			PublishVerificationResults: true,
			ProviderVersion:            "0.0.1",
			ProviderTags:               []string{"main"},
			StateHandlers:              messageStateHandlers(),
		},
		MessageHandlers: messageHandlers(),
	})

	assert.NoError(t, err)
}

func messageHandlers() message.MessageHandlers {
	return message.MessageHandlers{
		"a user1 registered event": func([]models.ProviderStateV3) (message.MessageBody, message.MessageMetadata, error) {
			event := eventSniffer.GetAndClearEvents()[0]
			return event, nil, nil
		},
		"a user1 details corrected event": func([]models.ProviderStateV3) (message.MessageBody, message.MessageMetadata, error) {
			event := eventSniffer.GetAndClearEvents()[0]
			return event, nil, nil
		},
	}
}

func messageStateHandlers() models.StateHandlers {
	return models.StateHandlers{
		"user1 has been registered": func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
			return nil, useCase.RegisterNewUser(context.Background(), testUser(1))
		},
		"user1 details have been corrected": func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
			return nil, useCase.CorrectUserDetails(context.Background(), testUser(1).Id, newTestUserDetails(1))
		},
	}
}

func emptyRepository() models.StateHandler {
	return func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
		return nil, repo.Clear(context.Background())
	}
}

func repositoryWith(users ...model.User) models.StateHandler {
	return func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
		ctx := context.Background()
		if err := repo.Clear(ctx); err != nil {
			return nil, err
		}

		for _, user := range users {
			if err := repo.AddUser(ctx, user); err != nil {
				return nil, err
			}
		}

		return nil, nil
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
