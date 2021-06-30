// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"github.com/go-pg/pg/v10/orm"
)

// Filter common filter interface.
type Filter interface {
	Apply(*orm.Query, condFn) *orm.Query
	ApplyGroup(*orm.Query, condGroupFn) *orm.Query
}

type condFn = func(condition string, params ...interface{}) *orm.Query

type condGroupFn = func(fn func(*orm.Query) (*orm.Query, error)) *orm.Query
