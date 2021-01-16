// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import (
	"encoding/json"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// Exists exists common filter.
type Exists struct {
	column string
	Value  bool `json:"value,omitempty"` // TODO: use pointer
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *Exists) UnmarshalJSON(b []byte) error {
	type alias Exists

	m1 := alias{}
	var m2 bool

	if err := json.Unmarshal(b, &m1); err == nil {
		f.Value = m1.Value
		return nil
	}

	if err := json.Unmarshal(b, &m2); err == nil {
		f.Value = m2
		return nil
	}

	return nil // TODO: return unsupported format error
}

// MarshalJSON custom JSON marshaler.
func (f *Exists) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

// NewExists initializes a new exists filter.
func NewExists(value bool) *Exists {
	return &Exists{
		Value: value,
	}
}

// Column set the column(s) for the exists filter.
func (f *Exists) Column(column string) *Exists {
	f.column = column
	return f
}

func (f *Exists) buildValue() string {
	if f.Value {
		return "IS NOT NULL"
	}
	return "IS NULL"
}

// Build build query.
func (f *Exists) Build(condFn condFn) *orm.Query {
	v := f.buildValue()
	return condFn("? ?", types.Ident(f.column), types.Safe(v))
}
