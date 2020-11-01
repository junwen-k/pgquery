// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package sorter_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/sorter"
	"github.com/stretchr/testify/assert"
)

type OrderTestItem struct {
	Id   int64
	Name string
	Age  int
}

func setupOrderTestItemTable(t *testing.T) {
	err := db.Model((*OrderTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &OrderTestItem{
			Name: fmt.Sprintf("name-%d", itemCount),
			Age:  itemCount,
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestSorterOrder(t *testing.T) {
	setupOrderTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With asc value":                  sorterOrderWithAscValue,
		"With desc value":                 sorterOrderWithDescValue,
		"With multiple columns asc value": sorterOrderWithMultipleColumnsAscValue,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func sorterOrderWithAscValue(t *testing.T) {
	var items []OrderTestItem
	q := db.Model(&items)

	sorter.NewOrderAsc().Column("age").Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 10) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
			assert.Equal(t, idx+1, item.Age)
		}
	}
}

func sorterOrderWithDescValue(t *testing.T) {
	var items []OrderTestItem
	q := db.Model(&items)

	sorter.NewOrderDesc().Column("age").Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 10) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", len(items)-idx), item.Name)
			assert.Equal(t, len(items)-idx, item.Age)
		}
	}
}

func sorterOrderWithMultipleColumnsAscValue(t *testing.T) {
	var items []OrderTestItem
	q := db.Model(&items)

	sorter.NewOrderDesc().Column("name").Build(q)
	sorter.NewOrderDesc().Column("age").Build(q)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 10) {
		for idx, item := range items {
			switch idx {
			case 8:
				assert.NotEmpty(t, item.Id)
				assert.Equal(t, fmt.Sprintf("name-%d", 10), item.Name)
				assert.Equal(t, 10, item.Age)
			case 9:
				assert.NotEmpty(t, item.Id)
				assert.Equal(t, fmt.Sprintf("name-%d", 1), item.Name)
				assert.Equal(t, 1, item.Age)
			default:
				assert.NotEmpty(t, item.Id)
				assert.Equal(t, fmt.Sprintf("name-%d", int(math.Abs(float64(idx-9)))), item.Name)
				assert.Equal(t, int(math.Abs(float64(idx-9))), item.Age)
			}
		}
	}
}
