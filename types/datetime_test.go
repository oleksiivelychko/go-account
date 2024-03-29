package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

const sampleDateTime = "2023-03-30 17:02:01"

type TestDateTimeDB struct {
	DateTime DateTime `json:"datetime"`
}

func TestTypes_MarshalDateTime(t *testing.T) {
	parsedTime, err := time.Parse(time.DateTime, sampleDateTime)
	if err != nil {
		t.Fatal(err)
	}

	marshalTo := &TestDateTimeDB{DateTime: DateTime(parsedTime)}
	marshaledJSON, err := json.Marshal(marshalTo)
	if err != nil {
		t.Fatal(err)
	}

	unmarshalTo := &TestDateTimeDB{}
	err = json.Unmarshal(marshaledJSON, &unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	datetime, err := unmarshalTo.DateTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Trim(string(datetime), "\"") != sampleDateTime {
		t.Errorf("datetime mismatch: %s != %s", strings.Trim(string(datetime), "\""), sampleDateTime)
	}
}

func TestTypes_UnmarshalDateTime(t *testing.T) {
	unmarshalTo := &TestDateTimeDB{}
	stringJSON := []byte(fmt.Sprintf(`{"datetime":"%s"}`, sampleDateTime))

	err := json.Unmarshal(stringJSON, &unmarshalTo)
	if err != nil {
		t.Fatal(err)
	}

	parsedTime, err := time.Parse(time.DateTime, sampleDateTime)
	if err != nil {
		t.Fatal(err)
	}

	datetimeDB := &TestDateTimeDB{DateTime: DateTime(parsedTime)}
	datetime, err := datetimeDB.DateTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Trim(string(datetime), "\"") != sampleDateTime {
		t.Errorf("datetime mismatch: %s != %s", strings.Trim(string(datetime), "\""), sampleDateTime)
	}
}
