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

type DatetimeRangeTestItem struct {
	Id        int64
	Name      string
	CreatedAt time.Time
}

func setupDatetimeRangeTestItemTable(t *testing.T) {
	err := db.Model((*DatetimeRangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &DatetimeRangeTestItem{
			Name:      fmt.Sprintf("name-%d", itemCount),
			CreatedAt: testTime.Add(time.Duration(itemCount) * time.Hour),
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterDatetimeRange(t *testing.T) {
	setupDatetimeRangeTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With after filter":   filterDatetimeRangeWithAfterFilter,
		"With before filter":  filterDatetimeRangeWithBeforeFilter,
		"With from filter":    filterDatetimeRangeWithFromFilter,
		"With to filter":      filterDatetimeRangeWithToFilter,
		"With complex filter": filterDatetimeRangeWithComplexFilter,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterDatetimeRangeWithAfterFilter(t *testing.T) {
	var items []DatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewDateTimeRange().After(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+6), item.Name)
			assert.Equal(t, testTime.Add(time.Duration(idx+6)*time.Hour).Year(), item.CreatedAt.Year())
			assert.Equal(t, testTime.Add(time.Duration(idx+6)*time.Hour).Month(), item.CreatedAt.Month())
			assert.Equal(t, testTime.Add(time.Duration(idx+6)*time.Hour).Day(), item.CreatedAt.Day())
			assert.Equal(t, testTime.Add(time.Duration(idx+6)*time.Hour).Second(), item.CreatedAt.Second())
		}
	}
}

func filterDatetimeRangeWithBeforeFilter(t *testing.T) {
	var items []DatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewDateTimeRange().Before(testTime.Add(6 * time.Hour)).Column("created_at").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Year(), item.CreatedAt.Year())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Month(), item.CreatedAt.Month())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Day(), item.CreatedAt.Day())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Second(), item.CreatedAt.Second())
		}
	}
}

func filterDatetimeRangeWithFromFilter(t *testing.T) {
	var items []DatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewDateTimeRange().From(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 6) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+5), item.Name)
			assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Year(), item.CreatedAt.Year())
			assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Month(), item.CreatedAt.Month())
			assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Day(), item.CreatedAt.Day())
			assert.Equal(t, testTime.Add(time.Duration(idx+5)*time.Hour).Second(), item.CreatedAt.Second())
		}
	}
}

func filterDatetimeRangeWithToFilter(t *testing.T) {
	var items []DatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewDateTimeRange().To(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Year(), item.CreatedAt.Year())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Month(), item.CreatedAt.Month())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Day(), item.CreatedAt.Day())
			assert.Equal(t, testTime.Add(time.Duration(idx+1)*time.Hour).Second(), item.CreatedAt.Second())
		}
	}
}

func filterDatetimeRangeWithComplexFilter(t *testing.T) {
	var items []DatetimeRangeTestItem
	q := db.Model(&items)

	filter.NewDateTimeRange().From(testTime.Add(2 * time.Hour)).To(testTime.Add(6 * time.Hour)).Column("created_at").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+2), item.Name)
			assert.Equal(t, testTime.Add(time.Duration(idx+2)*time.Hour).Year(), item.CreatedAt.Year())
			assert.Equal(t, testTime.Add(time.Duration(idx+2)*time.Hour).Month(), item.CreatedAt.Month())
			assert.Equal(t, testTime.Add(time.Duration(idx+2)*time.Hour).Day(), item.CreatedAt.Day())
			assert.Equal(t, testTime.Add(time.Duration(idx+2)*time.Hour).Second(), item.CreatedAt.Second())
		}
	}
}
