package bookmark

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	Field string
	Value any
	Err   error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("bookmark not found with %s = %v: %s", e.Field, e.Value, e.Err.Error())
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

type URLExistsError struct {
	URL string
	Err error
}

func (e *URLExistsError) Error() string {
	return fmt.Sprintf("a bookmark with URL %s already exists: %s", e.URL, e.Err.Error())
}

func (e *URLExistsError) Unwrap() error {
	return e.Err
}

func IsURLExists(err error) bool {
	var u *URLExistsError
	return errors.As(err, &u)
}
