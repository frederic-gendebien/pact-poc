package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/frederic-gendebien/pact-poc/application/server/internal/interfaces/http"
	"github.com/frederic-gendebien/pact-poc/application/server/pkg/domain/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/stretchr/testify/assert"
	"log"
	gohttp "net/http"
	"reflect"
	"testing"
)

var (
	pact *consumer.HTTPMockProviderV3
)

func init() {
	log.Println("setup pact environment")

	var err error
	pact, err = consumer.NewV3Pact(consumer.MockHTTPProviderConfig{
		Consumer: "user-client",
		Provider: "user-server-http",
		LogDir:   "../../../../tests/pact/logs",
		PactDir:  "../../../../tests/pact/pacts",
	})
	if err != nil {
		log.Fatalln("could not configure client pact: ", err)
	}

	log.Println("pact environment setup done")
}

func TestClientPact_RegisterNewUser(t *testing.T) {
	t.Run("Register New User", func(t *testing.T) {
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 does not exist",
				Parameters: userParameters(user1),
			}).
			UponReceiving("A new user registration request for user1").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodPut,
				Path:    matchers.Term("/users", "^/users$"),
				Headers: requestHeadersWithBody(),
				Body:    user1,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusCreated,
				Headers: responseHeadersWithoutBody(),
				Body:    nil,
			})

		verify(t, func(userClient *Client) error {
			return userClient.RegisterNewUser(context.Background(), user1)
		})
	})
	t.Run("Register An Existing User", func(t *testing.T) {
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 exists already",
				Parameters: userParameters(user1),
			}).
			UponReceiving("A new user registration request for user1").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodPut,
				Path:    matchers.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    user1,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusBadRequest,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user email : email1 already exists"),
			})

		verify(t, func(userClient *Client) error {
			if err := userClient.RegisterNewUser(context.Background(), user1); err == nil || !errors.Is(err, model.BadRequestError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.BadRequestError{}, err)
			}

			return nil
		})
	})
}

func TestClientPact_CorrectUserDetails(t *testing.T) {
	t.Run("Correct Unknown User Details", func(t *testing.T) {
		user1 := testUser(1)
		newUser1Details := newTestUserDetails(1)

		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 does not exist",
				Parameters: nil,
			}).
			UponReceiving("A correct user details request for user1").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodPut,
				Path:    matchers.Term("/users/user1/details", "^/users/[a-z0-9-]+/details$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    newUser1Details,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusNotFound,
				Headers: responseHeadersWithoutBody(),
				Body:    errorResponse("user with id: user1 was not found"),
			})

		verify(t, func(userClient *Client) error {
			if err := userClient.CorrectUserDetails(context.Background(), user1.Id, newUser1Details); err == nil || !errors.Is(err, model.NotFoundError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.NotFoundError{}, err)
			}

			return nil
		})
	})
	t.Run("Correct Existing User Details", func(t *testing.T) {
		user1 := testUser(1)
		newUser1Details := newTestUserDetails(1)

		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 exists already",
				Parameters: userParameters(user1),
			}).
			UponReceiving("A correct user details request for user1").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodPut,
				Path:    matchers.Term("/users/user1/details", "^/users/[a-z0-9-]+/details$"),
				Query:   nil,
				Headers: requestHeadersWithBody(),
				Body:    newUser1Details,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusAccepted,
				Headers: responseHeadersWithoutBody(),
				Body:    nil,
			})

		verify(t, func(userClient *Client) error {
			return userClient.CorrectUserDetails(context.Background(), user1.Id, newUser1Details)
		})
	})
}

func TestClientPact_DeleteUser(t *testing.T) {
	t.Run("Delete An Existing User", func(t *testing.T) {
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 exists",
				Parameters: userParameters(user1),
			}).
			UponReceiving("A delete user1 request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodDelete,
				Path:    matchers.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusAccepted,
				Headers: responseHeadersWithoutBody(),
				Body:    nil,
			})

		verify(t, func(userClient *Client) error {
			return userClient.DeleteUser(context.Background(), testUser(1).Id)
		})
	})
	t.Run("Delete An Unknown User", func(t *testing.T) {
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 does not exist",
				Parameters: nil,
			}).
			UponReceiving("A delete user1 request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodDelete,
				Path:    matchers.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusNotFound,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user with id: user1 was not found"),
			})

		verify(t, func(userClient *Client) error {
			if err := userClient.DeleteUser(context.Background(), user1.Id); err == nil || !errors.Is(err, model.NotFoundError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.NotFoundError{}, err)
			}

			return nil
		})
	})
}

func TestClientPact_ListAllUsers(t *testing.T) {
	t.Run("List All Users When There Are Many", func(t *testing.T) {
		expectedUsers := []model.User{
			testUser(1),
			testUser(2),
			testUser(3),
			testUser(4),
			testUser(5),
		}
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "Many users exist",
				Parameters: usersParameters(expectedUsers),
			}).
			UponReceiving("A list all users request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodGet,
				Path:    matchers.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
				Body:    expectedUsers,
			})

		verify(t, func(userClient *Client) error {
			next := make(chan bool)
			defer close(next)

			users, err := userClient.ListAllUsers(context.Background(), next)
			if err != nil {
				return fmt.Errorf("could not list all users: %v", err)
			}

			userList := make([]model.User, 0, 5)
			for user := range users {
				userList = append(userList, user)
				next <- true
			}

			if !reflect.DeepEqual(userList, expectedUsers) {
				return fmt.Errorf("expected: %v, but got %v", expectedUsers, userList)
			}

			return nil
		})
	})
	t.Run("List All Users When There Are None", func(t *testing.T) {
		expectedUsers := make([]model.User, 0)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "No users exist",
				Parameters: nil,
			}).
			UponReceiving("A list all users request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodGet,
				Path:    matchers.Term("/users", "^/users$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
				Body:    expectedUsers,
			})

		verify(t, func(userClient *Client) error {
			next := make(chan bool)
			defer close(next)

			users, err := userClient.ListAllUsers(context.Background(), next)
			if err != nil {
				return fmt.Errorf("could not list all users: %v", err)
			}

			userList := make([]model.User, 0, 5)
			for user := range users {
				userList = append(userList, user)
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
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 exists",
				Parameters: userParameters(user1),
			}).
			UponReceiving("A find user1 by id request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodGet,
				Path:    matchers.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusOK,
				Headers: responseHeadersWithBody(),
				Body:    user1,
			})

		verify(t, func(userClient *Client) error {
			user, err := userClient.FindUserById(context.Background(), user1.Id)
			if err != nil {
				return fmt.Errorf("could not find user by id: %v", err)
			}

			if !reflect.DeepEqual(user, user1) {
				return fmt.Errorf("expected: %v, but got %v", user1, user)
			}

			return nil
		})
	})
	t.Run("Find An Unknown User By Id", func(t *testing.T) {
		user1 := testUser(1)
		pact.AddInteraction().
			Given(models.ProviderStateV3{
				Name:       "The user1 does not exist",
				Parameters: nil,
			}).
			UponReceiving("A find user1 by id request").
			WithCompleteRequest(consumer.Request{
				Method:  gohttp.MethodGet,
				Path:    matchers.Term("/users/user1", "^/users/[a-z0-9-]+$"),
				Query:   nil,
				Headers: requestHeadersWithoutBody(),
				Body:    nil,
			}).
			WithCompleteResponse(consumer.Response{
				Status:  gohttp.StatusNotFound,
				Headers: responseHeadersWithBody(),
				Body:    errorResponse("user with id: user1 was not found"),
			})

		verify(t, func(userClient *Client) error {
			if _, err := userClient.FindUserById(context.Background(), user1.Id); err == nil || !errors.Is(err, model.NotFoundError{}) {
				return fmt.Errorf("a %v was expected, but found: %v", model.NotFoundError{}, err)
			}

			return nil
		})
	})
}

func requestHeadersWithBody() matchers.MapMatcher {
	return matchers.MapMatcher{
		"Accept":       matchers.Term("application/json; charset=utf-8", `application\/json`),
		"Content-Type": matchers.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func requestHeadersWithoutBody() matchers.MapMatcher {
	return matchers.MapMatcher{
		"Accept": matchers.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func responseHeadersWithBody() matchers.MapMatcher {
	return matchers.MapMatcher{
		"Content-Type": matchers.Term("application/json; charset=utf-8", `application\/json`),
	}
}

func responseHeadersWithoutBody() matchers.MapMatcher {
	return matchers.MapMatcher{}
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

func userParameters(user model.User) map[string]interface{} {
	return map[string]interface{}{
		"Id":           user.Id,
		"Email":        user.Email,
		"Details.Name": user.Details.Name,
	}
}

func usersParameters(users []model.User) map[string]interface{} {
	if users == nil {
		return nil
	}

	result := make(map[string]interface{})
	for _, user := range users {
		result[string(user.Id)] = userParameters(user)
	}

	return result
}

func newTestUserDetails(number int) model.UserDetails {
	return model.UserDetails{
		Name: fmt.Sprintf("new_name%d", number),
	}
}

func errorResponse(message string) http.ErrorResponse {
	return http.ErrorResponse{
		Message: message,
	}
}

func verify(t *testing.T, verification func(client *Client) error) {
	assert.NoError(t,
		pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			return verification(NewClient(fmt.Sprintf("http://%s:%s", config.Host, config.Port)))
		}),
	)
}
