package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type DateTime time.Time

//goland:noinspection GoMixedReceiverTypes
func (datetime *DateTime) MarshalJSON() ([]byte, error) {
	instantDatetime := time.Time(*datetime)
	return []byte(fmt.Sprintf("\"%v\"", instantDatetime.Format(time.DateTime))), nil
}

//goland:noinspection GoMixedReceiverTypes
func (datetime *DateTime) UnmarshalJSON(b []byte) (err error) {
	parsedTime, err := time.Parse(time.DateTime, strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	*datetime = DateTime(parsedTime)
	return
}

//goland:noinspection GoMixedReceiverTypes
func (datetime DateTime) Value() (driver.Value, error) {
	var instantTime time.Time
	instantDatetime := time.Time(datetime)
	if instantDatetime.UnixNano() == instantTime.UnixNano() {
		return nil, nil
	}

	return instantDatetime, nil
}

//goland:noinspection GoMixedReceiverTypes
func (datetime *DateTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*datetime = DateTime(value)
		return nil
	}

	return fmt.Errorf("unable to assign '%v' from DateTime to time.Time", v)
}
