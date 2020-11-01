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

type KeywordSearchTestItem struct {
	Id     int64
	Name   string
	Emails []string `pg:",array"`
}

func setupKeywordSearchTestItemTable(t *testing.T) {
	err := db.Model((*KeywordSearchTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	assert.NoError(t, err)

	for itemCount := 1; itemCount <= 10; itemCount++ {
		item := &KeywordSearchTestItem{
			Name:   fmt.Sprintf("name-%d", itemCount),
			Emails: make([]string, 0),
		}
		for emailCount := 1; emailCount <= 5; emailCount++ {
			item.Emails = append(item.Emails, fmt.Sprintf("email-%d(%d)@root", itemCount, emailCount))
		}
		_, err = db.Model(item).Insert()
		assert.NoError(t, err)
	}
}

func TestFilterKeywordSearch(t *testing.T) {
	setupKeywordSearchTestItemTable(t)

	tests := map[string]func(t *testing.T){
		"With default search":          filterKeywordSearchWithDefaultSearch,
		"With multiple columns search": filterKeywordSearchWithMultipleColumnsSearch,
		"With sensitive search":        filterKeywordSearchWithMatchAllSearch,
		"With sensitive start search":  filterKeywordSearchWithMatchStartSearch,
		"With sensitive end search":    filterKeywordSearchWithMatchEndSearch,
		"With case insensitive search": filterKeywordSearchWithCaseInsensitiveSearch,
		"With array search":            filterKeywordSearchWithArraySearch,
	}
	for name, test := range tests {
		t.Run(name, test)
	}
}

func filterKeywordSearchWithDefaultSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("name-1").Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 2) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
		assert.Equal(t, []string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}, items[0].Emails)

		assert.NotEmpty(t, items[1].Id)
		assert.Equal(t, "name-10", items[1].Name)
		assert.Equal(t, []string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}, items[1].Emails)
	}
}

func filterKeywordSearchWithMultipleColumnsSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("(1)@root").Column("name").Build(q.Where)
	filter.NewKeywordSearch("(1)@root").Column("emails,array").Build(q.WhereOr)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 10) {
		for idx, item := range items {
			assert.NotEmpty(t, item.Id)
			assert.Equal(t, fmt.Sprintf("name-%d", idx+1), item.Name)
			emails := make([]string, 0)
			for emailCount := 1; emailCount <= 5; emailCount++ {
				emails = append(emails, fmt.Sprintf("email-%d(%d)@root", idx+1, emailCount))
			}
			assert.Equal(t, emails, item.Emails)
		}
	}
}

func filterKeywordSearchWithMatchAllSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("name-1").MatchAll().Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 1) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
		assert.Equal(t, []string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}, items[0].Emails)
	}
}

func filterKeywordSearchWithMatchStartSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("name-1").MatchStart().Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 2) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
		assert.Equal(t, []string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}, items[0].Emails)

		assert.NotEmpty(t, items[1].Id)
		assert.Equal(t, "name-10", items[1].Name)
		assert.Equal(t, []string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}, items[1].Emails)
	}
}

func filterKeywordSearchWithMatchEndSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("-1").MatchEnd().Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 1) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
		assert.Equal(t, []string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}, items[0].Emails)
	}
}

func filterKeywordSearchWithCaseInsensitiveSearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("NAME-10").CaseInsensitive().Column("name").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 1) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-10", items[0].Name)
		assert.Equal(t, []string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}, items[0].Emails)
	}
}

func filterKeywordSearchWithArraySearch(t *testing.T) {
	var items []KeywordSearchTestItem
	q := db.Model(&items)

	filter.NewKeywordSearch("email-1").Column("emails,array").Build(q.Where)

	err := q.Select()
	assert.NoError(t, err)

	if assert.Len(t, items, 2) {
		assert.NotEmpty(t, items[0].Id)
		assert.Equal(t, "name-1", items[0].Name)
		assert.Equal(t, []string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}, items[0].Emails)

		assert.NotEmpty(t, items[1].Id)
		assert.Equal(t, "name-10", items[1].Name)
		assert.Equal(t, []string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}, items[1].Emails)
	}
}
