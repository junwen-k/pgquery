// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/filter"
	"github.com/stretchr/testify/assert"
)

type RangeTestItem struct {
	Id     int64
	Age    int
	Height int
}

func setupRangeTestItemTable(t *testing.T) {
	err := db.Model((*RangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &RangeTestItem{
			Age:    itemCount,
			Height: (itemCount * 10) + 130,
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterRange(t *testing.T) {
	setupRangeTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With greater than filter":       filterRangeWithGtFilter,
		"With greater than equal filter": filterRangeWithGteFilter,
		"With less than filter":          filterRangeWithLtFilter,
		"With less than equal filter":    filterRangeWithLteFilter,
		"With complex filter":            filterRangeWithComplexFilter,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterRangeWithGtFilter(t *testing.T) {
	var items []RangeTestItem
	q := db.Model(&items)

	filter.NewRange().GreaterThan(5).Column("age").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, idx+6, item.Age)
		}
	}
}

func filterRangeWithGteFilter(t *testing.T) {
	var items []RangeTestItem
	q := db.Model(&items)

	filter.NewRange().GreaterThanEqual(5).Column("age").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 6) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, idx+5, item.Age)
		}
	}
}

func filterRangeWithLtFilter(t *testing.T) {
	var items []RangeTestItem
	q := db.Model(&items)

	filter.NewRange().LessThan(5).Column("age").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 4) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, idx+1, item.Age)
		}
	}
}

func filterRangeWithLteFilter(t *testing.T) {
	var items []RangeTestItem
	q := db.Model(&items)

	filter.NewRange().LessThanEqual(5).Column("age").Build(q.WhereGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, idx+1, item.Age)
		}
	}
}

func filterRangeWithComplexFilter(t *testing.T) {
	var items []RangeTestItem
	q := db.Model(&items)

	filter.NewRange().GreaterThan(5).LessThan(8).Column("age").Build(q.WhereGroup)
	filter.NewRange().GreaterThan(5).LessThan(8).Column("height").Build(q.WhereOrGroup)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 2) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, idx+6, item.Age)
		}
	}
}
