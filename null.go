package fest

import (
	"database/sql"
	"encoding/json"
)

// NullInt64 is a wraper for sql.NullInt64 that works with JSON Unmarshal
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON wraps the json.Marshal function
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON wraps the json.Unmarshal function
func (v *NullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}

	return nil
}

// NullString is a wraper for sql.NullString that works with JSON Unmarshal
type NullString struct {
	sql.NullString
}

// MarshalJSON wraps the json.Marshal function
func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON wraps the json.Unmarshal function
func (v *NullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}

	return nil
}
