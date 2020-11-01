// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/filter"
	"github.com/stretchr/testify/assert"
)

type OffsetPaginationTestItem struct {
	Id   int64
	Name string
}

func setupOffsetPaginationTestItemTable(t *testing.T) {
	err := db.Model((*OffsetPaginationTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &OffsetPaginationTestItem{
			Name: fmt.Sprintf("name-%d", itemCount),
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterPagination(t *testing.T) {
	setupOffsetPaginationTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With page and limit": filterOffsetPaginationWithPageAndLimit,
		"With default limit":  filterOffsetPaginationWithDefaultLimit,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterOffsetPaginationWithPageAndLimit(t *testing.T) {
	var items []OffsetPaginationTestItem
	q := db.Model(&items)

	filter.NewOffsetPagination(1, 5).Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
		}
	}
}

func filterPaginationWithNoPage(t *testing.T) {
	var items []OffsetPaginationTestItem
	q := db.Model(&items)

	filter.NewOffsetPagination(0, 5).Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
		}
	}
}

func filterOffsetPaginationWithDefaultLimit(t *testing.T) {
	var items []OffsetPaginationTestItem
	q := db.Model(&items)

	filter.NewOffsetPagination(1, 0).DefaultLimit(5).Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
		}
	}
}

func filterPaginationWithNoLimit(t *testing.T) {
	var items []OffsetPaginationTestItem
	q := db.Model(&items)

	filter.NewOffsetPagination(1, 0).Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 10) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
		}
	}
}
