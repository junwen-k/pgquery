// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter

import "github.com/go-pg/pg/v10/orm"

// TODO: add typed misconfigured query error
// TODO: add tests for JSON marshal / unmarshal

// Filter common filter interface.
type Filter interface {
	Build()
}

type condFn = func(condition string, params ...interface{}) *orm.Query

type condGroupFn = func(fn func(*orm.Query) (*orm.Query, error)) *orm.Query
