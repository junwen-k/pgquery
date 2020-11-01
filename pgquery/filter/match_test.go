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

type MatchTestItem struct {
	Id   int64
	Name string
}

func setupMatchTestItemTable(t *testing.T) {
	err := db.Model((*MatchTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &MatchTestItem{
			Name: fmt.Sprintf("name-%d", itemCount),
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterMatch(t *testing.T) {
	setupMatchTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With single value":    filterMatchWithSingleValue,
		"With multiple values": filterMatchWithMultipleValues,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterMatchWithSingleValue(t *testing.T) {
	var items []MatchTestItem
	q := db.Model(&items)

	filter.NewMatch("name-1").Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 1) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
	}
}

func filterMatchWithMultipleValues(t *testing.T) {
	var items []MatchTestItem
	q := db.Model(&items)

	filter.NewMatch("name-1", "name-2", "name-3").Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 3) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
		}
	}
}
