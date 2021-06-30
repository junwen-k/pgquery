// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"encoding/json"
	"errors"

	"github.com/go-pg/pg/v10/types"
)

// Match match common filter.
type Match struct {
	column string
	Values []interface{} `json:"values,omitempty"`
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *Match) UnmarshalJSON(b []byte) error {
	type Alias Match

	m1 := Alias{}
	m2 := struct {
		*Alias
		Values interface{} `json:"values,omitempty"`
	}{Alias: (*Alias)(f)}
	m3 := make([]interface{}, 0)
	var m4 interface{}

	if err := json.Unmarshal(b, &m1); err == nil {
		f.Values = m1.Values
		return nil
	}

	if err := json.Unmarshal(b, &m2); err == nil {
		f.Values = []interface{}{m2.Values}
		return nil
	}

	if err := json.Unmarshal(b, &m3); err == nil {
		f.Values = m3
		return nil
	}

	if err := json.Unmarshal(b, &m4); err == nil {
		f.Values = []interface{}{m4}
		return nil
	}

	return errors.New("[Match]: unsupported format when unmarshalling json")
}

// MarshalJSON custom JSON marshaler.
func (f *Match) MarshalJSON() ([]byte, error) {
	switch {
	case len(f.Values) == 1:
		return json.Marshal(f.Values[0])
	default:
		return json.Marshal(f.Values)
	}
}

// NewMatch initializes a new match filter.
func NewMatch(column string) *Match {
	return &Match{
		column: column,
	}
}

// Column sets the column for the match filter.
func (f *Match) Column(column string) *Match {
	f.column = column
	return f
}

// Matches set value(s).
func (f *Match) Matches(values ...interface{}) *Match {
	f.Values = append(f.Values, values...)
	return f
}

// Appender returns parameters for cond appender.
func (f *Match) Appender() (string, interface{}, interface{}) {
	switch {
	case len(f.Values) > 1:
		return "? IN (?)", types.Ident(f.column), types.In(f.Values)
	default:
		return "? = ?", types.Ident(f.column), f.Values[0]
	}
}
