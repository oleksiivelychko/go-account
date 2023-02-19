package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type GormDateTime time.Time

//goland:noinspection GoMixedReceiverTypes
func (dt *GormDateTime) MarshalJSON() ([]byte, error) {
	timestamp := time.Time(*dt)
	return []byte(fmt.Sprintf("\"%v\"", timestamp.Format(time.DateTime))), nil
}

//goland:noinspection GoMixedReceiverTypes
func (dt *GormDateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*dt = GormDateTime(date)
	return
}

//goland:noinspection GoMixedReceiverTypes
func (dt GormDateTime) Value() (driver.Value, error) {
	var zeroTimestamp time.Time
	timestamp := time.Time(dt)
	if timestamp.UnixNano() == zeroTimestamp.UnixNano() {
		return nil, nil
	}
	return timestamp, nil
}

//goland:noinspection GoMixedReceiverTypes
func (dt *GormDateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*dt = GormDateTime(value)
		return nil
	}
	return fmt.Errorf("unable to convert %v to timestamp", v)
}
