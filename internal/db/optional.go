package db

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/cmessinides/mnemonic/internal/tag"
)

type Optional[T driver.Value] struct {
	HasValue bool
	Value    T
}

type OptionalString Optional[string]

func (o *OptionalString) UnmarshalParam(param string) error {
	o.HasValue = true
	o.Value = param
	return nil
}

type OptionalTags Optional[tag.Tags]

func (o *OptionalTags) UnmarshalParams(params []string) error {
	o.HasValue = true
	o.Value = tag.Tags(params)
	return nil
}

type OptionalBool Optional[bool]

func (o *OptionalBool) UnmarshalParam(param string) error {
	o.HasValue = true
	if param == "true" || param == "1" {
		o.Value = true
	}
	return nil
}

type OptionalSetter[T driver.Value] struct {
	Optional[T]
	Column string
}

func (o OptionalSetter[T]) UpdateQuery(query io.StringWriter, args []any) ([]any, error) {
	if !o.HasValue {
		return args, nil
	}

	_, err := query.WriteString(fmt.Sprintf("%s = ?, ", o.Column))
	if err != nil {
		return args, err
	}

	return append(args, o.Value), nil
}
