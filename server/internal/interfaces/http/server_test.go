package http

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/domain"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/infrastructure/persistence/inmemory"
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/usecase"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	dir, _  = os.Getwd()
	pactDir = fmt.Sprintf("%s/../../../tests/pact", dir)
	logDir  = fmt.Sprintf("%s/../../../tests/pact/log", dir)
	pact    = dsl.Pact{
		Provider:                 "pactpocServer",
		LogDir:                   logDir,
		PactDir:                  pactDir,
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}
)

func TestServer_AddUser(t *testing.T) {
	port := startServerWith(emptyRepository())
	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:            fmt.Sprintf("http://localhost:%d", port),
		PactURLs:                   nil,
		BrokerURL:                  "",
		ConsumerVersionSelectors:   nil,
		Tags:                       nil,
		ProviderTags:               nil,
		ProviderBranch:             "",
		ProviderStatesSetupURL:     "",
		Provider:                   "",
		BrokerUsername:             "",
		BrokerPassword:             "",
		BrokerToken:                "",
		FailIfNoPactsFound:         false,
		PublishVerificationResults: false,
		ProviderVersion:            "",
		CustomProviderHeaders:      nil,
		StateHandlers:              nil,
		BeforeEach:                 nil,
		AfterEach:                  nil,
		RequestFilter:              nil,
		CustomTLSConfig:            nil,
		EnablePending:              false,
		IncludeWIPPactsSince:       nil,
		PactLogDir:                 "",
		PactLogLevel:               "",
		Verbose:                    false,
		Args:                       nil,
	})
}

func emptyRepository() domain.UserRepository {
	return inmemory.NewUserRepository()
}

func startServerWith(repository domain.UserRepository) int {
	port, err := utils.GetFreePort()
	if err != nil {
		log.Fatalf("could not get free port: %v", err)
	}

	err = os.Setenv("PORT", strconv.Itoa(port))
	if err != nil {
		log.Fatalf("could not set port as environment variable: %v", err)
	}

	server := NewServer(usecase.NewDefaultUserCase(repository))
	go func() {
		log.Fatalf("server crashed: %v", server.Start())
	}()

	return port
}
