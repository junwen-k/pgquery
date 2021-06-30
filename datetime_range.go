// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// DateTimeRange common datetime range filter.
type DateTimeRange struct {
	column           string
	layouts          []string
	gtMarshalLayout  string
	gteMarshalLayout string
	ltMarshalLayout  string
	lteMarshalLayout string
	Gt               *time.Time `json:"after,omitempty"`
	Gte              *time.Time `json:"from,omitempty"`
	Lt               *time.Time `json:"before,omitempty"`
	Lte              *time.Time `json:"to,omitempty"`
}

// MarshalJSON custom JSON marshaler.
func (f *DateTimeRange) MarshalJSON() ([]byte, error) {
	type alias DateTimeRange

	m1 := struct {
		Gt  string `json:"after,omitempty"`
		Gte string `json:"from,omitempty"`
		Lt  string `json:"before,omitempty"`
		Lte string `json:"to,omitempty"`
		*alias
	}{alias: (*alias)(f)}

	if f.Gt != nil {
		if f.gtMarshalLayout == "" {
			return nil, errors.New("[DateTimeRange]: gtMarshalLayout is not specified for marshal json")
		}
		m1.Gt = f.Gt.Format(f.gtMarshalLayout)
	}
	if f.Lt != nil {
		if f.ltMarshalLayout == "" {
			return nil, errors.New("[DateTimeRange]: ltMarshalLayout is not specified for marshal json")
		}
		m1.Lt = f.Lt.Format(f.ltMarshalLayout)
	}
	if f.Gte != nil {
		if f.gteMarshalLayout == "" {
			return nil, errors.New("[DateTimeRange]: gteMarshalLayout is not specified for marshal json")
		}
		m1.Gte = f.Gte.Format(f.gteMarshalLayout)
	}
	if f.Lte != nil {
		if f.lteMarshalLayout == "" {
			return nil, errors.New("[DateTimeRange]: lteMarshalLayout is not specified for marshal json")
		}
		m1.Lte = f.Lte.Format(f.lteMarshalLayout)
	}

	return json.Marshal(m1)
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *DateTimeRange) UnmarshalJSON(b []byte) error {
	type alias DateTimeRange

	m1 := struct {
		Gt  string `json:"after,omitempty"`
		Gte string `json:"from,omitempty"`
		Lt  string `json:"before,omitempty"`
		Lte string `json:"to,omitempty"`
		*alias
	}{alias: (*alias)(f)}

	if len(f.layouts) <= 0 {
		return errors.New("[DateTimeRange]: layouts are not specified for unmarshal json")
	}

	err := json.Unmarshal(b, &m1)
	if err != nil {
		return errors.New("[DateTimeRange]: unsupported format when unmarshalling json")
	}

	for _, layout := range f.layouts {
		if f.Gt == nil && m1.Gt != "" {
			after, err := time.Parse(layout, m1.Gt)
			if err != nil {
				continue
			}
			f.gtMarshalLayout = layout
			f.Gt = &after
		}
		if f.Lt == nil && m1.Lt != "" {
			before, err := time.Parse(layout, m1.Lt)
			if err != nil {
				continue
			}
			f.ltMarshalLayout = layout
			f.Lt = &before
		}
		if f.Gte == nil && m1.Gte != "" {
			from, err := time.Parse(layout, m1.Gte)
			if err != nil {
				continue
			}
			f.gteMarshalLayout = layout
			f.Gte = &from
		}
		if f.Lte == nil && m1.Lte != "" {
			to, err := time.Parse(layout, m1.Lte)
			if err != nil {
				continue
			}
			f.lteMarshalLayout = layout
			f.Lte = &to
		}
	}

	return nil
}

// NewDateTimeRange initializes a new datetime range filter.
func NewDateTimeRange(column string, layouts ...string) *DateTimeRange {
	return &DateTimeRange{
		column:           column,
		layouts:          append(layouts, time.RFC3339),
		gtMarshalLayout:  time.RFC3339,
		gteMarshalLayout: time.RFC3339,
		ltMarshalLayout:  time.RFC3339,
		lteMarshalLayout: time.RFC3339,
	}
}

// Column sets the column for the datetime range filter.
func (f *DateTimeRange) Column(column string) *DateTimeRange {
	f.column = column
	return f
}

// Layout sets the parsing layout(s) for the datetime range filter.
func (f *DateTimeRange) Layout(layouts ...string) *DateTimeRange {
	f.layouts = append(f.layouts, layouts...)
	return f
}

// AfterMarshalLayout set marshal layout for after (gt).
func (f *DateTimeRange) AfterMarshalLayout(layout string) *DateTimeRange {
	f.gtMarshalLayout = layout
	return f
}

// BeforeMarshalLayout set marshal layout for before (lt).
func (f *DateTimeRange) BeforeMarshalLayout(layout string) *DateTimeRange {
	f.ltMarshalLayout = layout
	return f
}

// FromMarshalLayout set marshal layout for from (gte).
func (f *DateTimeRange) FromMarshalLayout(layout string) *DateTimeRange {
	f.gteMarshalLayout = layout
	return f
}

// ToMarshalLayout set marshal layout for to (lte).
func (f *DateTimeRange) ToMarshalLayout(layout string) *DateTimeRange {
	f.lteMarshalLayout = layout
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

// Appender returns parameters for cond group appender.
func (f *DateTimeRange) Appender() applyFn {
	return func(q *orm.Query) (*orm.Query, error) {
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
	}
}
