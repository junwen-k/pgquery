// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/filter"
	"github.com/stretchr/testify/assert"
)

type RelativeDatetimeRangeTestItem struct {
	Id        int64
	Name      string
	CreatedAt time.Time `pg:"type:timestamp"` // Without timezone.
}

func setupRelativeDatetimeRangeTestItemTable(t *testing.T) {
	err := db.Model((*RelativeDatetimeRangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &RelativeDatetimeRangeTestItem{
			Name: fmt.Sprintf("name-%d", itemCount),
		}
		if itemCount%2 == 0 {
			item.CreatedAt = testTime.Add(time.Duration(itemCount) * time.Hour)
		} else {
			item.CreatedAt = testTime.Add(-time.Duration(itemCount) * time.Hour)
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterRelativeDatetimeRange(t *testing.T) {
	setupRelativeDatetimeRangeTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With hour ago filter":            filterRelativeDatetimeRangeWithHourAgoFilter,
		"With hour upcoming filter":       filterRelativeDatetimeRangeWithHourUpcomingFilter,
		"With hour ago / upcoming filter": filterRelativeDatetimeRangeWithHourAgoUpcomingFilter,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterRelativeDatetimeRangeWithHourAgoFilter(t *testing.T) {
	var items []RelativeDatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewRelativeDateTimeRange(testTime).Column("created_at").AgoHour(5).Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 3) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			switch idx {
			case 0:
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Second(), item.CreatedAt.Second())
			case 1:
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Second(), item.CreatedAt.Second())
			case 2:
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Second(), item.CreatedAt.Second())
			}
		}
	}
}

func filterRelativeDatetimeRangeWithHourUpcomingFilter(t *testing.T) {
	var items []RelativeDatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewRelativeDateTimeRange(testTime).Column("created_at").UpcomingHour(5).Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 2) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			switch idx {
			case 0:
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Second(), item.CreatedAt.Second())
			case 1:
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Second(), item.CreatedAt.Second())
			}
		}
	}
}

func filterRelativeDatetimeRangeWithHourAgoUpcomingFilter(t *testing.T) {
	var items []RelativeDatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewRelativeDateTimeRange(testTime).Column("created_at").AgoHour(5).UpcomingHour(5).Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			switch idx {
			case 0:
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+2)*time.Hour).Second(), item.CreatedAt.Second())
			case 1:
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(time.Duration(idx+3)*time.Hour).Second(), item.CreatedAt.Second())
			case 2:
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+4)*time.Hour).Second(), item.CreatedAt.Second())
			case 3:
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Second(), item.CreatedAt.Second())
			case 4:
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Year(), item.CreatedAt.Year())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Month(), item.CreatedAt.Month())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Day(), item.CreatedAt.Day())
				assert.Equal(t, testTime.Add(-time.Duration(idx+6)*time.Hour).Second(), item.CreatedAt.Second())
			}
		}
	}
}
