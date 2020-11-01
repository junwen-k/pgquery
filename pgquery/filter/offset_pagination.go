// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import "github.com/go-pg/pg/v10/orm"

// OffsetPagination pagination common filter. Offset pagination is skipped when no limit is provided.
type OffsetPagination struct {
	defaultLimit int
	Page         int `json:"page,omitempty"`
	Limit        int `json:"limit,omitempty"`
}

// NewOffsetPagination initializes a new pagination filter.
func NewOffsetPagination(page, limit int) *OffsetPagination {
	return &OffsetPagination{
		Page:  page,
		Limit: limit,
	}
}

func (f *OffsetPagination) init() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = f.defaultLimit
	}
}

// DefaultLimit sets the default limit for the pagination filter.
func (f *OffsetPagination) DefaultLimit(limit int) *OffsetPagination {
	f.defaultLimit = limit
	return f
}

// Build build query.
func (f *OffsetPagination) Build(q *orm.Query) *orm.Query {
	f.init()
	if f.Limit > 0 {
		q.Limit(f.Limit)
		q.Offset((f.Page - 1) * f.Limit)
	}
	return q
}
