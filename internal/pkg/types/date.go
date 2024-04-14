package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

func (d *Date) Scan(src any) error {
	switch t := src.(type) {
	case time.Time:
		*d = Date(t)
	case *time.Time:
		if t == nil {
			return nil
		}
		*d = Date(*t)
	case string:
		v, err := ParseDate(t)
		if err != nil {
			return err
		}
		*d = v
	default:
		return fmt.Errorf("scan: unknown value for Date: %+v", src)
	}

	return nil
}

func (d *Date) Time() time.Time {
	if d == nil {
		return time.Time{}
	}

	return time.Time(*d)
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}

func ParseDate(val string) (Date, error) {
	t, err := time.Parse(time.DateOnly, val)
	if err != nil {
		return Date{}, err
	}

	return Date(t), nil
}
