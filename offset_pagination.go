// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import "github.com/go-pg/pg/v10/orm"

// OffsetPagination pagination common filter. Offset pagination is skipped when no limit is provided.
type OffsetPagination struct {
	Page  int  `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// NewOffsetPagination initializes a new pagination filter.
func NewOffsetPagination() *OffsetPagination {
	return &OffsetPagination{}
}

func (f *OffsetPagination) init() {
	if f.Page <= 0 {
		f.Page = 1
	}
}

// Offset sets offset value for the pagination filter.
func (f *OffsetPagination) Offset(page int, limit int) *OffsetPagination {
	f.Page = page
	f.Limit = &limit
	return f
}

// Build build query.
func (f *OffsetPagination) Build(q *orm.Query) *orm.Query {
	f.init()
	if limit := f.Limit; limit != nil {
		q.Limit(*limit)
		q.Offset((f.Page - 1) * *limit)
	}
	return q
}
