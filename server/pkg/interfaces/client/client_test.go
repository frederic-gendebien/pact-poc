package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/frederic-gendebien/poc-pact/server/internal/interfaces/http"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"github.com/pact-foundation/pact-go/dsl"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	pact       dsl.Pact
	userClient *Client
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		pact, userClient = setup()
		defer tearDown()

		exitCode := m.Run()

		if err := pact.WritePact(); err != nil {
			log.Fatalln("could not write pact: ", err)
		}

		return exitCode
	}())
}

func TestClientPact_RegisterNewUser(t *testing.T) {
	t.Run("Register New User", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 does not exist").
			UponReceiving("A new user registration for user1").
			WithRequest(dsl.Request{
				Method:  "PUT",
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    testUser(1),
			}).
			WillRespondWith(dsl.Response{
				Status:  201,
				Headers: responseHeaders(),
				Body:    nil,
			})

		verify(t, pact, func() error {
			return userClient.RegisterNewUser(context.Background(), testUser(1))
		})
	})
	t.Run("Register An Existing User", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user2 exists already").
			UponReceiving("A new user registration for user2").
			WithRequest(dsl.Request{
				Method:  "PUT",
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    testUser(2),
			}).
			WillRespondWith(dsl.Response{
				Status:  400,
				Headers: responseHeaders(),
				Body:    errorResponse("user2 exists already"),
			})

		verify(t, pact, func() error {
			if err := userClient.RegisterNewUser(context.Background(), testUser(2)); err == nil || !errors.Is(err, model.BadRequestError{}) {
				return fmt.Errorf("a bad request was expected, but found: %v", err)
			}

			return nil
		})
	})
}

func TestClientPact_FindUserById(t *testing.T) {
	t.Run("Find An Existing User By Id", func(t *testing.T) {
		expectedUser := testUser(1)
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 exists").
			UponReceiving("A find user1 by id request").
			WithRequest(dsl.Request{
				Method:  "GET",
				Path:    dsl.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  200,
				Headers: responseHeaders(),
				Body:    expectedUser,
			})

		verify(t, pact, func() error {
			user, err := userClient.FindUserById(context.Background(), expectedUser.GetId())
			if err != nil {
				return fmt.Errorf("could not find user by id: %v", err)
			}

			if !reflect.DeepEqual(user, NewUserFrom(expectedUser)) {
				return fmt.Errorf("expected: %v, but got %v", expectedUser, user)
			}

			return nil
		})
	})
}

func setup() (dsl.Pact, *Client) {
	log.Println("clearing pact folders")
	if err := os.RemoveAll("../../../../test/pact"); err != nil {
		log.Fatalln("could not clear pact folders: ", err)
	}
	log.Println("pact folders clearing done")

	log.Println("setup pact environment")
	pact := dsl.Pact{
		Consumer:                 "user-client",
		Provider:                 "user-server",
		LogDir:                   "../../../../test/pact/logs",
		PactDir:                  "../../../../test/pact/pacts",
		LogLevel:                 "INFO",
		DisableToolValidityCheck: true,
	}
	pact.Setup(true)
	userClient := NewClient(fmt.Sprintf("http://localhost:%d", pact.Server.Port))

	log.Println("pact environment setup done")
	return pact, userClient
}

func tearDown() {
	log.Println("tearing down pact environment")
	pact.Teardown()
	_ = userClient.Close()
	log.Println("pact environment tear down done")
}

func requestHeadersWithBody() dsl.MapMatcher {
	return dsl.MapMatcher{
		"Accept":       dsl.Term("application/json; charset=utf-8", `application\/json`),
		"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func requestHeadersWithoutBody() dsl.MapMatcher {
	return dsl.MapMatcher{
		"Accept": dsl.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func responseHeaders() dsl.MapMatcher {
	return dsl.MapMatcher{
		"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func testUser(number int) http.User {
	return http.User{
		Id:    fmt.Sprintf("user%d", number),
		Name:  fmt.Sprintf("name%d", number),
		Email: fmt.Sprintf("email%d", number),
	}
}

func errorResponse(message string) http.ErrorResponse {
	return http.ErrorResponse{
		Message: message,
	}
}

func verify(t *testing.T, pact dsl.Pact, verification func() error) {
	if err := pact.Verify(verification); err != nil {
		t.Fatalf("Error on verify: %v", err)
	}
}
