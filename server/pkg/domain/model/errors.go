package model

func NewBadRequest(message string) BadRequestError {
	return BadRequestError{Message: message}
}

type BadRequestError struct {
	Message string
}

func (b BadRequestError) Error() string {
	return b.Message
}

func (b BadRequestError) Is(err error) bool {
	_, ok := err.(BadRequestError)

	return ok
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

type NotFoundError struct {
	Message string
}

func (b NotFoundError) Error() string {
	return b.Message
}

func (b NotFoundError) Is(err error) bool {
	_, ok := err.(NotFoundError)

	return ok
}

func NewUnknownError(message string, err error) UnknownError {
	return UnknownError{
		Message: message,
		Err:     err,
	}
}

type UnknownError struct {
	Message string
	Err     error
}

func (b UnknownError) Error() string {
	return b.Message
}

func (b UnknownError) Is(err error) bool {
	_, ok := err.(UnknownError)

	return ok
}
