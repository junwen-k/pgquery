// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// RelativeDateTimeRangeUnitOption relative datetime range unit option.
type RelativeDateTimeRangeUnitOption struct {
	Century     *int `json:"century,omitempty"`
	Day         *int `json:"day,omitempty"`
	Decade      *int `json:"decade,omitempty"`
	Hour        *int `json:"hour,omitempty"`
	Microsecond *int `json:"microsecond,omitempty"`
	Millennium  *int `json:"millennium,omitempty"`
	Millisecond *int `json:"millisecond,omitempty"`
	Minute      *int `json:"minute,omitempty"`
	Month       *int `json:"month,omitempty"`
	Second      *int `json:"second,omitempty"`
	Week        *int `json:"week,omitempty"`
	Year        *int `json:"year,omitempty"`
}

func (o *RelativeDateTimeRangeUnitOption) buildValue(value int, plural, singular string) string {
	switch {
	case value > 1:
		return fmt.Sprintf("%d %s", value, plural)
	case value == 1:
		return fmt.Sprintf("%d %s", value, singular)
	default:
		return ""
	}
}

func (o *RelativeDateTimeRangeUnitOption) build() string {
	var units []string
	if o.Millennium != nil {
		if v := o.buildValue(*o.Millennium, "millenniums", "millennium"); v != "" {
			units = append(units, v)
		}
	}
	if o.Century != nil {
		if v := o.buildValue(*o.Century, "centuries", "century"); v != "" {
			units = append(units, v)
		}
	}
	if o.Decade != nil {
		if v := o.buildValue(*o.Decade, "decades", "decade"); v != "" {
			units = append(units, v)
		}
	}
	if o.Year != nil {
		if v := o.buildValue(*o.Year, "years", "year"); v != "" {
			units = append(units, v)
		}
	}
	if o.Month != nil {
		if v := o.buildValue(*o.Month, "months", "month"); v != "" {
			units = append(units, v)
		}
	}
	if o.Week != nil {
		if v := o.buildValue(*o.Week, "weeks", "week"); v != "" {
			units = append(units, v)
		}
	}
	if o.Day != nil {
		if v := o.buildValue(*o.Day, "days", "day"); v != "" {
			units = append(units, v)
		}
	}
	if o.Hour != nil {
		if v := o.buildValue(*o.Hour, "hours", "hour"); v != "" {
			units = append(units, v)
		}
	}
	if o.Minute != nil {
		if v := o.buildValue(*o.Minute, "minutes", "minute"); v != "" {
			units = append(units, v)
		}
	}
	if o.Second != nil {
		if v := o.buildValue(*o.Second, "seconds", "second"); v != "" {
			units = append(units, v)
		}
	}
	if o.Millisecond != nil {
		if v := o.buildValue(*o.Millisecond, "milliseconds", "millisecond"); v != "" {
			units = append(units, v)
		}
	}
	if o.Microsecond != nil {
		if v := o.buildValue(*o.Microsecond, "microseconds", "microsecond"); v != "" {
			units = append(units, v)
		}
	}
	return strings.Join(units, " ")
}

// RelativeDateTimeRange relative datetime range common filter.
type RelativeDateTimeRange struct {
	column        string
	layouts       []string
	marshalLayout string
	Ago           *RelativeDateTimeRangeUnitOption `json:"ago,omitempty"`
	Upcoming      *RelativeDateTimeRangeUnitOption `json:"upcoming,omitempty"`
	At            *time.Time                       `json:"at,omitempty"`
}

// MarshalJSON custom JSON marshaler.
func (f *RelativeDateTimeRange) MarshalJSON() ([]byte, error) {
	type alias RelativeDateTimeRange

	m1 := struct {
		At string `json:"at,omitempty"`
		*alias
	}{alias: (*alias)(f)}

	if f.At != nil {
		if f.marshalLayout == "" {
			return nil, errors.New("[RelativeDateTimeRange]: marshalLayout is not specified for marshal json")
		}
		m1.At = f.At.Format(f.marshalLayout)
	}

	return json.Marshal(m1)
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *RelativeDateTimeRange) UnmarshalJSON(b []byte) error {
	type alias RelativeDateTimeRange

	m1 := struct {
		At string `json:"at,omitempty"`
		*alias
	}{alias: (*alias)(f)}

	if len(f.layouts) <= 0 {
		return errors.New("[RelativeDateTimeRange]: layouts are not specified for unmarshal json")
	}

	err := json.Unmarshal(b, &m1)
	if err != nil {
		return errors.New("[RelativeDateTimeRange]: unsupported format when unmarshalling json")
	}

	for _, layout := range f.layouts {
		if f.At == nil && m1.At != "" {
			at, err := time.Parse(layout, m1.At)
			if err != nil {
				continue
			}
			f.marshalLayout = layout
			f.At = &at
		}
	}

	return nil
}

// NewRelativeDateTimeRange initializes a new relative datetime filter.
func NewRelativeDateTimeRange(column string, layouts ...string) *RelativeDateTimeRange {
	f := &RelativeDateTimeRange{
		column:        column,
		layouts:       append(layouts, time.RFC3339),
		marshalLayout: time.RFC3339,
	}
	f.init()
	return f
}

func (f *RelativeDateTimeRange) init() {
	if f.At == nil {
		now := time.Now()
		f.At = &now
	}
	if f.Ago == nil {
		f.Ago = &RelativeDateTimeRangeUnitOption{}
	}
	if f.Upcoming == nil {
		f.Upcoming = &RelativeDateTimeRangeUnitOption{}
	}
}

// Column sets the column for the relative datetime filter.
func (f *RelativeDateTimeRange) Column(column string) *RelativeDateTimeRange {
	f.column = column
	return f
}

// Layout sets the parsing layout(s) for the relative datetime filter.
func (f *RelativeDateTimeRange) Layout(layouts ...string) *RelativeDateTimeRange {
	f.layouts = append(f.layouts, layouts...)
	return f
}

// MarshalLayout set marshal layout.
func (f *RelativeDateTimeRange) MarshalLayout(layout string) *RelativeDateTimeRange {
	f.marshalLayout = layout
	return f
}

// AgoCentury set century for ago.
func (f *RelativeDateTimeRange) AgoCentury(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Century = &value
	return f
}

// AgoDay set day for ago.
func (f *RelativeDateTimeRange) AgoDay(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Day = &value
	return f
}

// AgoDecade set decade for ago.
func (f *RelativeDateTimeRange) AgoDecade(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Decade = &value
	return f
}

// AgoHour set hour for ago.
func (f *RelativeDateTimeRange) AgoHour(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Hour = &value
	return f
}

// AgoMicrosecond set microsecond for ago.
func (f *RelativeDateTimeRange) AgoMicrosecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Microsecond = &value
	return f
}

// AgoMillennium set millennium for ago.
func (f *RelativeDateTimeRange) AgoMillennium(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Millennium = &value
	return f
}

// AgoMillisecond set millisecond for ago.
func (f *RelativeDateTimeRange) AgoMillisecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Millisecond = &value
	return f
}

// AgoMinute set minute for ago.
func (f *RelativeDateTimeRange) AgoMinute(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Minute = &value
	return f
}

// AgoMonth set month for ago.
func (f *RelativeDateTimeRange) AgoMonth(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Month = &value
	return f
}

// AgoSecond set second for ago.
func (f *RelativeDateTimeRange) AgoSecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Second = &value
	return f
}

// AgoWeek set week for ago.
func (f *RelativeDateTimeRange) AgoWeek(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Week = &value
	return f
}

// AgoYear set year for ago.
func (f *RelativeDateTimeRange) AgoYear(value int) *RelativeDateTimeRange {
	f.init()
	f.Ago.Year = &value
	return f
}

// UpcomingCentury set century for upcoming.
func (f *RelativeDateTimeRange) UpcomingCentury(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Century = &value
	return f
}

// UpcomingDay set day for upcoming.
func (f *RelativeDateTimeRange) UpcomingDay(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Day = &value
	return f
}

// UpcomingDecade set decade for upcoming.
func (f *RelativeDateTimeRange) UpcomingDecade(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Decade = &value
	return f
}

// UpcomingHour set hour for upcoming.
func (f *RelativeDateTimeRange) UpcomingHour(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Hour = &value
	return f
}

// UpcomingMicrosecond set microsecond for upcoming.
func (f *RelativeDateTimeRange) UpcomingMicrosecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Microsecond = &value
	return f
}

// UpcomingMillennium set millennium for upcoming.
func (f *RelativeDateTimeRange) UpcomingMillennium(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Millennium = &value
	return f
}

// UpcomingMillisecond set millisecond for upcoming.
func (f *RelativeDateTimeRange) UpcomingMillisecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Millisecond = &value
	return f
}

// UpcomingMinute set minute for upcoming.
func (f *RelativeDateTimeRange) UpcomingMinute(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Minute = &value
	return f
}

// UpcomingMonth set month for upcoming.
func (f *RelativeDateTimeRange) UpcomingMonth(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Month = &value
	return f
}

// UpcomingSecond set second for upcoming.
func (f *RelativeDateTimeRange) UpcomingSecond(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Second = &value
	return f
}

// UpcomingWeek set week for upcoming.
func (f *RelativeDateTimeRange) UpcomingWeek(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Week = &value
	return f
}

// UpcomingYear set year for upcoming.
func (f *RelativeDateTimeRange) UpcomingYear(value int) *RelativeDateTimeRange {
	f.init()
	f.Upcoming.Year = &value
	return f
}

// AsAt set value.
func (f *RelativeDateTimeRange) AsAt(at time.Time) *RelativeDateTimeRange {
	f.At = &at
	return f
}

// Appender returns parameters for cond group appender.
func (f *RelativeDateTimeRange) Appender() applyFn {
	f.init()
	return func(q *orm.Query) (*orm.Query, error) {
		if ago := f.Ago.build(); ago != "" {
			q.Where("? >= ?::timestamp - interval ?", types.Ident(f.column), f.At.Format(time.RFC3339Nano), ago)
		} else {
			q.Where("? >= ?::timestamp", types.Ident(f.column), f.At.Format(time.RFC3339Nano))
		}
		if upcoming := f.Upcoming.build(); upcoming != "" {
			q.Where("? <= ?::timestamp + interval ?", types.Ident(f.column), f.At.Format(time.RFC3339Nano), upcoming)
		} else {
			q.Where("? <= ?::timestamp", types.Ident(f.column), f.At.Format(time.RFC3339Nano))
		}
		return q, nil
	}
}
