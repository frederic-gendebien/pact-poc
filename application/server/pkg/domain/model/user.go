package model

type UserId string

func (i UserId) IsInvalid() bool {
	return i == ""
}

type Email string

func (e Email) IsInvalid() bool {
	return e == ""
}

type User struct {
	Id      UserId      `json:"id"`
	Details UserDetails `json:"details"`
	Email   Email       `json:"email"`
}

func (u User) CorrectDetails(newDetails UserDetails) User {
	u.Details = newDetails
	return u
}

func (u User) IsInvalid() bool {
	return u.Id.IsInvalid() ||
		u.Details.IsInvalid() ||
		u.Email.IsInvalid()
}

func (u User) Invalid() error {
	if u.IsInvalid() {
		return NewBadRequest("user content is invalid")
	}

	return nil
}

func (u User) InvalidAfter(err error) error {
	if err != nil {
		return err
	}

	return u.Invalid()
}

type UserDetails struct {
	Name string `json:"name"`
}

func (d UserDetails) IsInvalid() bool {
	return d.Name == ""
}

func (d UserDetails) Invalid() error {
	if d.IsInvalid() {
		return NewBadRequest("user details is invalid")
	}

	return nil
}

func (d UserDetails) InvalidAfter(err error) error {
	if err != nil {
		return err
	}

	return d.Invalid()
}
