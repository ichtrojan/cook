package helpers

import (
	"time"
)

type JSONTime struct {
	time.Time
}

func (t JSONTime) Json() string {
	return t.Truncate(time.Second).Format("2006-01-02T15:04:05.000000Z")
}
