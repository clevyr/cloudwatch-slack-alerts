package alert

import (
	"errors"
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05.000-0700"

type Time struct {
	time.Time
}

var ErrTimeInvalidInput = errors.New("Time.UnmarshalJSON: input is not JSON")

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("%w: %q", ErrTimeInvalidInput, data)
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
