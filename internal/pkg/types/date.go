package types

import (
	"encoding/json"
	"time"
)

type Date time.Time

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
