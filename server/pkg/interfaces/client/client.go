package client

import (
	"context"
	"github.com/frederic-gendebien/poc-pact/server/pkg/domain/model"
	"github.com/go-resty/resty/v2"
)

func NewClient(url string) *Client {
	client := resty.New()
	client.SetBaseURL(url)
	client.SetHeader("Accept", "application/json; charset=utf-8")
	return &Client{
		client: client,
	}
}

type Client struct {
	client *resty.Client
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) RegisterNewUser(ctx context.Context, newUser model.User) error {
	response, err := c.client.R().
		SetBody(NewUserFrom(newUser)).
		Put("/users")

	if err != nil {
		return model.NewUnknownError("could not register new user", err)
	}

	_, err = bodyOrError(response, emptyBody())

	return err
}

func (c *Client) DeleteUser(ctx context.Context, userId string) error {
	response, err := c.client.R().
		SetPathParam("user_id", userId).
		Delete("/users/{user_id}")

	if err != nil {
		return model.NewUnknownError("could not delete user", err)
	}

	_, err = bodyOrError(response, emptyBody())

	return err
}

func (c *Client) ListAllUsers(ctx context.Context, next <-chan bool) (<-chan model.User, error) {
	response, err := c.client.R().
		Get("/users")

	if err != nil {
		return nil, model.NewUnknownError("could not list all users", err)
	}

	users, err := bodyOrError(response, usersProvider())
	if err != nil {
		return nil, err
	}

	results := make(chan model.User)
	go func() {
		defer close(results)

		for _, user := range users.([]User) {
			results <- user
			select {
			case needNext := <-next:
				if !needNext {
					break
				}
			}
		}
	}()

	return results, nil
}

func (c *Client) FindUserById(ctx context.Context, userId string) (model.User, error) {
	response, err := c.client.R().
		SetPathParam("user_id", userId).
		Get("/users/{user_id}")

	if err != nil {
		return nil, model.NewUnknownError("could not find user by id", err)
	}

	user, err := bodyOrError(response, userProvider())

	return user.(model.User), err
}
