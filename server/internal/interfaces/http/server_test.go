package http

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/infrastructure/persistence/inmemory"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/usecase"
	"context"
	"fmt"
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
	pactBrokerUrl   string
	pactBrokerToken string
	port            int
	pact            dsl.Pact
	repository      *inmemory.UserRepository
	useCase         usecase.UserUseCase
	server          *Server
)

func init() {
	pactBrokerUrl = os.Getenv(pactBrokerUrlPropertyName)
	if pactBrokerUrl == "" {
		log.Fatalln("missing environment property: ", pactBrokerUrlPropertyName)
	}

	pactBrokerToken = os.Getenv(pactBrokerTokenPropertyName)
	if pactBrokerUrl == "" {
		log.Fatalln("missing environment property: ", pactBrokerTokenPropertyName)
	}

	var err error
	port, err = utils.GetFreePort()
	if err != nil {
		log.Fatalf("could not find free port: %v", err)
	}

	if err := os.Setenv("PORT", strconv.Itoa(port)); err != nil {
		log.Fatalf("could not set port environment variable: %v", err)
	}

	pact = dsl.Pact{
		Provider:                 "user-server",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}

	repository = inmemory.NewUserRepository()
	useCase = usecase.NewUserUseCase(repository)
	server = NewServer(useCase)
	go func() {
		log.Println(server.Start())
	}()
}

func TestServerPact(t *testing.T) {
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:            fmt.Sprintf("http://127.0.0.1:%d", port),
		Tags:                       []string{"master"},
		BrokerURL:                  pactBrokerUrl,
		BrokerToken:                pactBrokerToken,
		FailIfNoPactsFound:         true,
		ProviderVersion:            "0.0.1",
		PublishVerificationResults: true,
		StateHandlers:              stateHandlers(),
	})
	if err != nil {
		t.Fatalf("server verifaction failed: %v", err)
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

func repositoryWith(users ...User) types.StateHandler {
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

func testUser(number int) User {
	return User{
		Id:    fmt.Sprintf("user%d", number),
		Name:  fmt.Sprintf("name%d", number),
		Email: fmt.Sprintf("email%d", number),
	}
}
