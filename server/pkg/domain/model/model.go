package model

type User interface {
	GetId() string
	GetName() string
	GetEmail() string
}
