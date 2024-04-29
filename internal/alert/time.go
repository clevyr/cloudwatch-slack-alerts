package alert

import (
	"errors"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05.000-0700"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	var err error
	t.Time, err = time.Parse(TimeFormat, string(data))
	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	val := t.Format(TimeFormat)
	return []byte(val), nil
}
