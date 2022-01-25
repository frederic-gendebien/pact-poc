package model

type UserId string
type Email string

type User struct {
	Id    UserId `json:"id"`
	Email Email  `json:"email"`
}
