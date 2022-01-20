package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/interfaces/http"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/pact-foundation/pact-go/dsl"
	"log"
	gohttp "net/http"
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
			UponReceiving("A new user registration request for user1").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodPut,
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    testUser(1),
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusCreated,
				Headers: responseHeadersWithoutBody(),
				Body:    nil,
			})

		verify(t, pact, func() error {
			return userClient.RegisterNewUser(context.Background(), testUser(1))
		})
	})
	t.Run("Register An Existing User", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 exists already").
			UponReceiving("A new user registration request for user1").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodPut,
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    testUser(1),
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusBadRequest,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user with id: user1 already exists"),
			})

		verify(t, pact, func() error {
			if err := userClient.RegisterNewUser(context.Background(), testUser(1)); err == nil || !errors.Is(err, model.BadRequestError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.BadRequestError{}, err)
			}

			return nil
		})
	})
}

func TestClientPact_ListAllUsers(t *testing.T) {
	t.Run("List All Users When There Are Many", func(t *testing.T) {
		pact.Interactions = nil
		expectedUsers := []User{
			testUser(1),
			testUser(2),
			testUser(3),
			testUser(4),
			testUser(5),
		}
		pact.AddInteraction().
			Given("Many users exist").
			UponReceiving("A list all users request").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodGet,
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
				Body:    expectedUsers,
			})

		verify(t, pact, func() error {
			next := make(chan bool)
			defer close(next)

			users, err := userClient.ListAllUsers(context.Background(), next)
			if err != nil {
				return fmt.Errorf("could not list all users: %v", err)
			}

			userList := make([]User, 0, 5)
			for user := range users {
				userList = append(userList, NewUserFrom(user))
				next <- true
			}

			if !reflect.DeepEqual(userList, expectedUsers) {
				return fmt.Errorf("expected: %v, but got %v", expectedUsers, userList)
			}

			return nil
		})
	})
	t.Run("List All Users When There Are None", func(t *testing.T) {
		pact.Interactions = nil
		expectedUsers := make([]User, 0)
		pact.AddInteraction().
			Given("No users exist").
			UponReceiving("A list all users request").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodGet,
				Path:    dsl.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
				Body:    expectedUsers,
			})

		verify(t, pact, func() error {
			next := make(chan bool)
			defer close(next)

			users, err := userClient.ListAllUsers(context.Background(), next)
			if err != nil {
				return fmt.Errorf("could not list all users: %v", err)
			}

			userList := make([]User, 0, 5)
			for user := range users {
				userList = append(userList, NewUserFrom(user))
				next <- true
			}

			if !reflect.DeepEqual(userList, expectedUsers) {
				return fmt.Errorf("expected: %v, but got %v", expectedUsers, userList)
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
				Method:  gohttp.MethodGet,
				Path:    dsl.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
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
	t.Run("Find An Unknown User By Id", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 does not exist").
			UponReceiving("A find user1 by id request").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodGet,
				Path:    dsl.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusNotFound,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user with id: user1 was not found"),
			})

		verify(t, pact, func() error {
			if _, err := userClient.FindUserById(context.Background(), testUser(1).GetId()); err == nil || !errors.Is(err, model.NotFoundError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.NotFoundError{}, err)
			}

			return nil
		})
	})
}

func TestClientPact_DeleteUser(t *testing.T) {
	t.Run("Delete An Existing User", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 exists").
			UponReceiving("A delete user1 request").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodDelete,
				Path:    dsl.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusAccepted,
				Headers: responseHeadersWithoutBody(),
				Body:    nil,
			})

		verify(t, pact, func() error {
			err := userClient.DeleteUser(context.Background(), testUser(1).GetId())
			if err != nil {
				return fmt.Errorf("could not delete user: %v", err)
			}

			return nil
		})
	})
	t.Run("Delete An Unknown User", func(t *testing.T) {
		pact.Interactions = nil
		pact.AddInteraction().
			Given("The user1 does not exist").
			UponReceiving("A delete user1 request").
			WithRequest(dsl.Request{
				Method:  gohttp.MethodDelete,
				Path:    dsl.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WillRespondWith(dsl.Response{
				Status:  gohttp.StatusNotFound,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user with id: user1 was not found"),
			})

		verify(t, pact, func() error {
			if err := userClient.DeleteUser(context.Background(), testUser(1).GetId()); err == nil || !errors.Is(err, model.NotFoundError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.NotFoundError{}, err)
			}

			return nil
		})
	})
}

func setup() (dsl.Pact, *Client) {
	log.Println("clearing pact folders")
	if err := os.RemoveAll("../../../../tests/pact"); err != nil {
		log.Fatalln("could not clear pact folders: ", err)
	}
	log.Println("pact folders clearing done")

	log.Println("setup pact environment")
	pact := dsl.Pact{
		Consumer:                 "user-client",
		Provider:                 "user-server",
		LogDir:                   "../../../../tests/pact/logs",
		PactDir:                  "../../../../tests/pact/pacts",
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

func responseHeadersWithBody() dsl.MapMatcher {
	return dsl.MapMatcher{
		"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func responseHeadersWithoutBody() dsl.MapMatcher {
	return dsl.MapMatcher{}
}

func testUser(number int) User {
	return User{
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
