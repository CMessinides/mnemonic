package tag

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tags []string

func (t *Tags) Scan(src any) error {
	var err error
	switch s := src.(type) {
	case []byte:
		err = json.Unmarshal(s, &t)
	case string:
		err = json.Unmarshal([]byte(s), &t)
	case nil:
		return nil
	default:
		err = fmt.Errorf("cannot handle value of type %T", s)
	}

	if err != nil {
		return fmt.Errorf("unable to parse tags: %w", err)
	}

	return nil
}

func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}
