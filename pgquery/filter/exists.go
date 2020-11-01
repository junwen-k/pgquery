// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import (
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// Exists exists common filter.
type Exists struct {
	column string
	Value  bool `json:"value,omitempty"`
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
