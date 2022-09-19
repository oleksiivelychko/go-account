package models

import (
	"fmt"
	"time"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("02-01-2006 15:04:05 MST"))
	return []byte(stamp), nil
}
