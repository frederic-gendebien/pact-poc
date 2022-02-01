package http

import (
	"context"
	"fmt"
	inmemorypers "github.com/frederic-gendebien/pact-poc/application/server/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/events"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/config/environment"
	"github.com/frederic-gendebien/pact-poc/lib/eventbus"
	inmemoryevb "github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"log"
	"os"
	"strconv"
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
	port            int
	repository      *inmemorypers.UserRepository
	eventBus        *inmemoryevb.EventBus
	eventSniffer    *eventbus.EventSniffer
	useCase         usecase.UserUseCase
	server          *Server
)

func init() {
	configuration = environment.NewConfiguration()
	pactBrokerUrl = configuration.GetStringOrCrash(pactBrokerUrlPropertyName)
	pactBrokerToken = configuration.GetStringOrCrash(pactBrokerTokenPropertyName)

	var err error
	port, err = utils.GetFreePort()
	if err != nil {
		log.Fatalf("could not find free port: %v", err)
	}

	if err := os.Setenv("PORT", strconv.Itoa(port)); err != nil {
		log.Fatalf("could not set port environment variable: %v", err)
	}

	repository = inmemorypers.NewUserRepository()
	eventBus = inmemoryevb.NewEventBus()
	eventSniffer = eventbus.NewEventSniffer(eventBus)
	useCase = usecase.NewUserUseCase(repository, eventBus)
	server = NewServer(useCase)

	go func() {
		log.Println(server.Start())
	}()
}

func TestServerHTTPPact(t *testing.T) {
	pact := dsl.Pact{
		Provider:                 "user-server-http",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}

	if _, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:            fmt.Sprintf("http://127.0.0.1:%d", port),
		Tags:                       []string{"main"},
		BrokerURL:                  pactBrokerUrl,
		BrokerToken:                pactBrokerToken,
		FailIfNoPactsFound:         true,
		ProviderVersion:            "0.0.1",
		ProviderTags:               []string{"main"},
		PublishVerificationResults: true,
		StateHandlers:              stateHandlers(),

		PactLogDir: "../../../../tests/pact/logs",
	}); err != nil {
		t.Fatalf("server http verifaction failed: %v", err)
	}
}

func TestServerMessagePact(t *testing.T) {
	eventSniffer.Clear()
	if err := eventSniffer.Listen(events.NewUserRegistered{}); err != nil {
		t.Fatal(err)
	}

	pact := dsl.Pact{
		Provider:                 "user-server-async",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
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
			return eventSniffer.GetEvents()[0], nil
		},
	}
}

func messageStateHandlers() dsl.StateHandlers {
	return dsl.StateHandlers{
		"user1 has been registered": func(state dsl.State) error {
			return useCase.RegisterNewUser(context.Background(), testUser(1))
		},
	}
}

func stateHandlers() types.StateHandlers {
	return types.StateHandlers{
		"The user1 does not exist": emptyRepository(),
		"The user1 exists":         repositoryWith(testUser(1)),
		"The user1 exists already": repositoryWith(testUser(1)),
		"Many users exist": repositoryWith(
			testUser(1),
			testUser(2),
			testUser(3),
			testUser(4),
			testUser(5),
		),
		"No users exist": emptyRepository(),
	}
}

func emptyRepository() types.StateHandler {
	return func() error {
		return repository.Clear(context.Background())
	}
}

func repositoryWith(users ...model.User) func() error {
	return func() error {
		ctx := context.Background()
		if err := repository.Clear(ctx); err != nil {
			return err
		}

		for _, user := range users {
			if err := repository.AddUser(ctx, user); err != nil {
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
