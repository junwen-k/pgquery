// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// DateTimeRange common datetime range filter.
type DateTimeRange struct {
	column string
	layout string
	Gt     *time.Time `json:"after,omitempty"`
	Gte    *time.Time `json:"from,omitempty"`
	Lt     *time.Time `json:"before,omitempty"`
	Lte    *time.Time `json:"to,omitempty"`
}

// MarshalJSON custom JSON marshaler.
func (f *DateTimeRange) MarshalJSON() ([]byte, error) {
	var t *struct {
		Gt  string `json:"after,omitempty"`
		Gte string `json:"from,omitempty"`
		Lt  string `json:"before,omitempty"`
		Lte string `json:"to,omitempty"`
	}
	if f.Gt != nil {
		t.Gt = f.Gt.Format(f.layout)
	}
	if f.Lt != nil {
		t.Lt = f.Lt.Format(f.layout)
	}
	if f.Gte != nil {
		t.Gte = f.Gte.Format(f.layout)
	}
	if f.Lte != nil {
		t.Lte = f.Lte.Format(f.layout)
	}
	return json.Marshal(t)
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *DateTimeRange) UnmarshalJSON(b []byte) error {
	var t *struct {
		Gt  string `json:"after,omitempty"`
		Gte string `json:"from,omitempty"`
		Lt  string `json:"before,omitempty"`
		Lte string `json:"to,omitempty"`
	}
	if err := json.Unmarshal(b, &t); err == nil {
		if t.Gt != "" {
			after, err := time.Parse(f.layout, t.Gt)
			if err == nil {
				f.Gt = &after
			}
		}
		if t.Lt != "" {
			before, err := time.Parse(f.layout, t.Lt)
			if err == nil {
				f.Lt = &before
			}
		}
		if t.Gte != "" {
			from, err := time.Parse(f.layout, t.Gte)
			if err == nil {
				f.Gte = &from
			}
		}
		if t.Lte != "" {
			to, err := time.Parse(f.layout, t.Lte)
			if err == nil {
				f.Lte = &to
			}
		}
	}
	return nil
}

// NewDateTimeRange initializes a new datetime range filter.
func NewDateTimeRange() *DateTimeRange {
	return &DateTimeRange{
		layout: time.RFC3339,
	}
}

// Column sets the column for the datetime range filter.
func (f *DateTimeRange) Column(column string) *DateTimeRange {
	f.column = column
	return f
}

// Layout sets the parsing layout for the datetime range filter.
func (f *DateTimeRange) Layout(layout string) *DateTimeRange {
	f.layout = layout
	return f
}

// After set value for after (gt).
func (f *DateTimeRange) After(value time.Time) *DateTimeRange {
	f.Gt = &value
	return f
}

// Before set value for before (lt).
func (f *DateTimeRange) Before(value time.Time) *DateTimeRange {
	f.Lt = &value
	return f
}

// From set value for from (gte).
func (f *DateTimeRange) From(value time.Time) *DateTimeRange {
	f.Gte = &value
	return f
}

// To set value for to (lte).
func (f *DateTimeRange) To(value time.Time) *DateTimeRange {
	f.Lte = &value
	return f
}

// Build build query.
func (f *DateTimeRange) Build(condGroupFn condGroupFn) *orm.Query {
	return condGroupFn(func(q *orm.Query) (*orm.Query, error) {
		if f.Lt != nil {
			q.Where("? < ?", types.Ident(f.column), f.Lt.Format(time.RFC3339Nano))
		}
		if f.Lte != nil {
			q.Where("? <= ?", types.Ident(f.column), f.Lte.Format(time.RFC3339Nano))
		}
		if f.Gte != nil {
			q.Where("? >= ?", types.Ident(f.column), f.Gte.Format(time.RFC3339Nano))
		}
		if f.Gt != nil {
			q.Where("? > ?", types.Ident(f.column), f.Gt.Format(time.RFC3339Nano))
		}
		return q, nil
	})
}
