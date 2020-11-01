// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import (
	"strings"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// KeywordSearch keyword search common filter.
type KeywordSearch struct {
	column          string
	caseInsensitive bool
	matchAll        bool
	matchStart      bool
	matchEnd        bool
	Value           string `json:"value,omitempty"`
}

// NewKeywordSearch initializes a new keyword search filter.
func NewKeywordSearch(value string) *KeywordSearch {
	return &KeywordSearch{
		Value: value,
	}
}

// Column set the column for the keyword search filter. Suffix column with ",array" to use array search.
func (f *KeywordSearch) Column(column string) *KeywordSearch {
	f.column = column
	return f
}

// CaseInsensitive set keyword case insensitive for the keyword search filter.
func (f *KeywordSearch) CaseInsensitive() *KeywordSearch {
	f.caseInsensitive = true
	return f
}

// MatchAll set keyword match all for the keyword search filter.
func (f *KeywordSearch) MatchAll() *KeywordSearch {
	f.matchAll = true
	return f
}

// MatchStart set keyword match start for the keyword search filter.
func (f *KeywordSearch) MatchStart() *KeywordSearch {
	f.matchStart = true
	return f
}

// MatchEnd set keyword match end for the keyword search filter.
func (f *KeywordSearch) MatchEnd() *KeywordSearch {
	f.matchEnd = true
	return f
}

func (f *KeywordSearch) buildValue() string {
	v := f.Value
	if f.matchAll {
		return v
	}
	if !f.matchStart {
		v = "%" + f.Value
	}
	if !f.matchEnd {
		v += "%"
	}
	return v
}

func (f *KeywordSearch) buildLike() string {
	if f.caseInsensitive {
		return "ILIKE"
	}
	return "LIKE"
}

func (f *KeywordSearch) buildColumn(column string) interface{} {
	if strings.HasSuffix(column, ",array") {
		return orm.SafeQuery("array_to_string(?, ?)", types.Ident(strings.TrimSuffix(column, ",array")), ",")
	}
	return types.Ident(column)
}

// Build build query.
func (f *KeywordSearch) Build(condFn condFn) *orm.Query {
	v := f.buildValue()
	column := f.buildColumn(f.column)
	like := f.buildLike()
	return condFn("? ? ?", column, types.Safe(like), v)
}
