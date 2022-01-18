package client

import (
	"bitbucket.org/fredericgendebien/pact-poc/server/internal/interfaces/http"
	"bitbucket.org/fredericgendebien/pact-poc/server/pkg/domain/model"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	gohttp "net/http"
)

func bodyOrError(response *resty.Response, provider func([]byte) (interface{}, error)) (interface{}, error) {
	switch response.RawResponse.StatusCode {
	case gohttp.StatusOK:
		return provider(response.Body())
	case gohttp.StatusCreated:
		return provider(response.Body())
	case gohttp.StatusAccepted:
		return provider(response.Body())
	case gohttp.StatusNoContent:
		return nil, nil
	case gohttp.StatusBadRequest:
		return nil, model.NewBadRequest(errorMessage(response))
	case gohttp.StatusNotFound:
		return nil, model.NewNotFoundError(errorMessage(response))
	default:
		return nil, model.NewUnknownError(errorMessage(response), nil)
	}
}

func emptyBody() func([]byte) (interface{}, error) {
	return func(bytes []byte) (interface{}, error) {
		return nil, nil
	}
}

func userProvider() func([]byte) (interface{}, error) {
	return func(bytes []byte) (interface{}, error) {
		user := User{}
		if err := json.Unmarshal(bytes, &user); err != nil {
			return user, err
		}

		return user, nil
	}
}

func usersProvider() func([]byte) (interface{}, error) {
	return func(bytes []byte) (interface{}, error) {
		var users []User
		if err := json.Unmarshal(bytes, &users); err != nil {
			return nil, err
		}

		return users, nil
	}
}

func errorMessage(response *resty.Response) string {
	errorResponse := http.ErrorResponse{}
	if err := json.Unmarshal(response.Body(), &errorResponse); err != nil {
		return fmt.Sprintf("could not read response message: %v", err)
	}

	return errorResponse.Message
}
