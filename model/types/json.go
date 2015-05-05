package types

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONObject represents a JSON object.
type JSONObject map[string]interface{}

// Scan implements the sql.Scanner interface.
func (j *JSONObject) Scan(val interface{}) error {
	if b, ok := val.([]byte); ok {
		return json.Unmarshal(b, j)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (j JSONObject) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}

	return json.Marshal(j)
}

// JSONArray represents a JSONArray.
type JSONArray []interface{}

// Scan implements the sql.Scanner interface.
func (j *JSONArray) Scan(val interface{}) error {
	if b, ok := val.([]byte); ok {
		return json.Unmarshal(b, j)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}

	return json.Marshal(j)
}
