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

type ExistsTestItem struct {
	Id   int64
	Name string
}

func setupExistsTestItemTable(t *testing.T) {
	err := db.Model((*ExistsTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		name := ""
		if itemCount%2 == 0 {
			name = fmt.Sprintf("name-%d", itemCount)
		}
		item := &ExistsTestItem{
			Name: name,
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterExists(t *testing.T) {
	setupExistsTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With true value":  filterExistsWithTrueValue,
		"With false value": filterExistsWithFalseValue,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterExistsWithTrueValue(t *testing.T) {
	var items []ExistsTestItem
	q := db.Model(&items)

	filter.NewExists(true).Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", (idx+1)*2), item.Name)
		}
	}
}

func filterExistsWithFalseValue(t *testing.T) {
	var items []ExistsTestItem
	q := db.Model(&items)

	filter.NewExists(false).Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 5) {
		for _, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Empty(t, item.Name)
		}
	}
}
