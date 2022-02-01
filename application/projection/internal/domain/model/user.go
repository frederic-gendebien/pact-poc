package model

type UserId string
type Email string

type User struct {
	Id    UserId `json:"id"`
	Name  string `json:"name"`
	Email Email  `json:"email"`
}

func (u User) UpdateWith(user User) User {
	return User{
		Id:    u.Id,
		Name:  valueOr(user.Name, u.Name),
		Email: Email(valueOr(string(user.Email), string(u.Email))),
	}
}

func valueOr(value string, defaultValue string) string {
	if value != "" {
		return value
	}

	return defaultValue
}
