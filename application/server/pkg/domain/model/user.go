package model

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) IsInvalid() bool {
	return u.Id == "" ||
		u.Name == "" ||
		u.Email == ""
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
