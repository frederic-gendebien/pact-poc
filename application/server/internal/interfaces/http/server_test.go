package http

import (
	"context"
	"fmt"
	inmemorypers "github.com/frederic-gendebien/pact-poc/application/server/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/usecase"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"github.com/frederic-gendebien/pact-poc/lib/config/environment"
	inmemoryevb "github.com/frederic-gendebien/pact-poc/lib/eventbus/inmemory"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/pact-foundation/pact-go/v2/utils"
	"github.com/stretchr/testify/assert"
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
	useCase = usecase.NewUserUseCase(repository, eventBus)
	server = NewServer(useCase)

	go func() {
		log.Println(server.Start())
	}()
}

func TestServerHTTPPact(t *testing.T) {
	pact := &provider.HTTPVerifier{}
	err := pact.VerifyProvider(t, provider.VerifyRequest{
		BrokerURL:                  pactBrokerUrl,
		BrokerToken:                pactBrokerToken,
		Provider:                   "user-server-http",
		ProviderBaseURL:            fmt.Sprintf("http://127.0.0.1:%d", port),
		ProviderVersion:            "0.0.1",
		ProviderTags:               []string{"main"},
		Tags:                       []string{"main"},
		FailIfNoPactsFound:         true,
		PublishVerificationResults: true,
		StateHandlers:              stateHandlers(),
		PactDirs:                   []string{"../../../../tests/pact/pacts"},

		//LogDir:                   "../../../../tests/pact/logs",
		//DisableToolValidityCheck: true,
		//LogLevel:                 "INFO",
		//PactLogDir:               "../../../../tests/pact/logs",
	})

	assert.NoError(t, err)
}

func stateHandlers() models.StateHandlers {
	return models.StateHandlers{
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

func emptyRepository() models.StateHandler {
	return func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
		return nil, repository.Clear(context.Background())
	}
}

func repositoryWith(users ...model.User) models.StateHandler {
	return func(setup bool, state models.ProviderStateV3) (models.ProviderStateV3Response, error) {
		ctx := context.Background()
		if err := repository.Clear(ctx); err != nil {
			return nil, err
		}

		for _, user := range users {
			if err := repository.AddUser(ctx, user); err != nil {
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
