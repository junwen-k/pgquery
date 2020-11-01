// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import (
	"encoding/json"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// Match match common filter.
type Match struct {
	column string
	Values []interface{} `json:"values,omitempty"`
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *Match) UnmarshalJSON(b []byte) error {
	var t *struct {
		Values string `json:"values,omitempty"`
	}
	if err := json.Unmarshal(b, &t); err == nil {
		f.Values = []interface{}{t.Values}
	}
	var tt *struct {
		Values []interface{} `json:"values,omitempty"`
	}
	if err := json.Unmarshal(b, &tt); err == nil {
		f.Values = tt.Values
	}
	return nil
}

// NewMatch initializes a new match filter.
func NewMatch(values ...interface{}) *Match {
	return &Match{Values: values}
}

// Column sets the column for the match filter.
func (f *Match) Column(column string) *Match {
	f.column = column
	return f
}

// Build build query.
func (f *Match) Build(condFn condFn) *orm.Query {
	switch {
	case len(f.Values) > 1:
		return condFn("? IN (?)", types.Ident(f.column), types.In(f.Values))
	case len(f.Values) == 1:
		return condFn("? = ?", types.Ident(f.column), f.Values[0])
	default:
		return nil
	}
}
