package bookmark

import (
	"errors"
	"fmt"
	"reflect"
)

type NotFoundError struct {
	Field string
	Value any
	Err   error
}

func (e *NotFoundError) Error() string {
	op := "="
	if reflect.TypeOf(e.Value).Kind() == reflect.Slice {
		op = "in"
	}

	msg := fmt.Sprintf("no bookmark found where %s %s %v", e.Field, op, e.Value)
	if e.Err != nil {
		msg = msg + ": " + e.Err.Error()
	}

	return msg
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

func IsNotFound(err error) bool {
	var n *NotFoundError
	return errors.As(err, &n)
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
