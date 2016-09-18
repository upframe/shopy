package models

import (
	"database/sql"
	"encoding/json"
)

// NullInt64JSON is a wraper for sql.NullInt64 that works with JSON Unmarshal
type NullInt64JSON struct {
	sql.NullInt64
}

// MarshalJSON wraps the json.Marshal function
func (v NullInt64JSON) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON wraps the json.Unmarshal function
func (v *NullInt64JSON) UnmarshalJSON(data []byte) error {
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

// NullStringJSON is a wraper for sql.NullString that works with JSON Unmarshal
type NullStringJSON struct {
	sql.NullString
}

// MarshalJSON wraps the json.Marshal function
func (v NullStringJSON) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON wraps the json.Unmarshal function
func (v *NullStringJSON) UnmarshalJSON(data []byte) error {
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
