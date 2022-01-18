package main

import (
	"context"
	"github.com/frederic-gendebien/poc-pact/server/pkg/interfaces/client"
	"github.com/google/uuid"
	"log"
)

var (
	url   string
	users *client.Client
)

func init() {
	users = client.NewClient("http://localhost:8080")
}

func main() {
	//createUsers()
	listUsers()
}

func createUsers() {
	for i := 0; i < 15; i++ {
		userId := uuid.New().String()
		if err := users.RegisterNewUser(context.Background(), client.User{
			Id:    userId,
			Name:  "Frederic Gendebien",
			Email: "frederic.gendebien@gmail.com",
		}); err != nil {
			log.Fatalln(err)
		}
	}
}

func listUsers() {
	next := make(chan bool)
	done := make(chan bool)
	defer close(done)

	go func() {
		defer close(next)

		userz, err := users.ListAllUsers(context.Background(), next)
		if err != nil {
			log.Fatalln(err)
		}

		for user := range userz {
			log.Println(user)
			next <- true
		}

		done <- true
	}()

	<-done
	log.Println("done")
}