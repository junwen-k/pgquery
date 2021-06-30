// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"encoding/json"
	"errors"

	"github.com/go-pg/pg/v10/types"
)

// NOTE: rename to NotNull
// TODO: if nil pointer, will not apply filter

// Exists exists common filter.
type Exists struct {
	column string
	Value  *bool `json:"value,omitempty"`
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *Exists) UnmarshalJSON(b []byte) error {
	type alias Exists

	m1 := alias{}
	var m2 *bool

	if err := json.Unmarshal(b, &m1); err == nil {
		f.Value = m1.Value
		return nil
	}

	if err := json.Unmarshal(b, &m2); err == nil {
		f.Value = m2
		return nil
	}

	return errors.New("[Exists]: unsupported format when unmarshalling json")
}

// MarshalJSON custom JSON marshaler.
func (f *Exists) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

// NewExists initializes a new exists filter.
func NewExists(column string) *Exists {
	return &Exists{
		column: column,
	}
}

// Column set the column(s) for the exists filter.
func (f *Exists) Column(column string) *Exists {
	f.column = column
	return f
}

// Exists set value.
func (f *Exists) Exists(value bool) *Exists {
	f.Value = &value
	return f
}

// ShouldExists set value to true.
func (f *Exists) ShouldExists() *Exists {
	return f.Exists(true)
}

// ShouldNotExists set value to false.
func (f *Exists) ShouldNotExists() *Exists {
	return f.Exists(false)
}

func (f *Exists) buildValue() string {
	if f.Value != nil && *f.Value {
		return "IS NOT NULL"
	}
	return "IS NULL"
}

// Appender returns parameters for cond appender.
func (f *Exists) Appender() (string, interface{}, interface{}) {
	v := f.buildValue()
	return "? ?", types.Ident(f.column), types.Safe(v)
}
