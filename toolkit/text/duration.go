package text

import (
	"fmt"
	"strings"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	text := strings.Trim(string(b), "\"")
	if text == "null" || text == "" {
		return nil
	}
	duration, err := time.ParseDuration(text)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Duration(*d).String())), nil
}

func (d *Duration) Duration() time.Duration {
	return time.Duration(*d)
}
