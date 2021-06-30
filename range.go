// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// Range range common filter.
type Range struct {
	column string
	Gt     *int `json:"gt,omitempty"`
	Gte    *int `json:"gte,omitempty"`
	Lt     *int `json:"lt,omitempty"`
	Lte    *int `json:"lte,omitempty"`
}

// NewRange initializes a new range filter.
func NewRange(column string) *Range {
	return &Range{
		column: column,
	}
}

// Column set the column for the range filter.
func (f *Range) Column(column string) *Range {
	f.column = column
	return f
}

// GreaterThan set value for greater than (gt).
func (f *Range) GreaterThan(value int) *Range {
	f.Gt = &value
	return f
}

// GreaterThanEqual set value for greater than (gte).
func (f *Range) GreaterThanEqual(value int) *Range {
	f.Gte = &value
	return f
}

// LessThan set value for less than (lt).
func (f *Range) LessThan(value int) *Range {
	f.Lt = &value
	return f
}

// LessThanEqual set value for less than equal (lte).
func (f *Range) LessThanEqual(value int) *Range {
	f.Lte = &value
	return f
}

// Build build query.
func (f *Range) Build(condGroupFn condGroupFn) *orm.Query {
	return condGroupFn(func(q *orm.Query) (*orm.Query, error) {
		if f.Lt != nil {
			q.Where("? < ?", types.Ident(f.column), f.Lt)
		}
		if f.Lte != nil {
			q.Where("? <= ?", types.Ident(f.column), f.Lte)
		}
		if f.Gte != nil {
			q.Where("? >= ?", types.Ident(f.column), f.Gte)
		}
		if f.Gt != nil {
			q.Where("? > ?", types.Ident(f.column), f.Gt)
		}
		return q, nil
	})
}
