// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"encoding/json"
	"errors"
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
	Value           *string `json:"value,omitempty"`
}

// UnmarshalJSON custom JSON unmarshaler.
func (f *KeywordSearch) UnmarshalJSON(b []byte) error {
	type alias KeywordSearch

	m1 := alias{}
	var m2 *string

	if err := json.Unmarshal(b, &m1); err == nil {
		f.Value = m1.Value
		return nil
	}

	if err := json.Unmarshal(b, &m2); err == nil {
		f.Value = m2
		return nil
	}

	return errors.New("[KeywordSearch]: unsupported format when unmarshalling json")
}

// MarshalJSON custom JSON marshaler.
func (f *KeywordSearch) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

// NewKeywordSearch initializes a new keyword search filter.
func NewKeywordSearch(column string) *KeywordSearch {
	return &KeywordSearch{
		column: column,
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

// Keyword set value.
func (f *KeywordSearch) Keyword(keyword string) *KeywordSearch {
	f.Value = &keyword
	return f
}

func (f *KeywordSearch) buildValue() string {
	var v string
	if f.Value != nil {
		v = *f.Value
	}
	if f.matchAll {
		return v
	}
	if !f.matchStart {
		v = "%" + v
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

// Appender returns parameters for cond appender.
func (f *KeywordSearch) Appender() (string, interface{}, interface{}, interface{}) {
	v := f.buildValue()
	column := f.buildColumn(f.column)
	like := f.buildLike()
	return "? ? ?", column, types.Safe(like), v
}
