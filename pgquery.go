// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery

import (
	"github.com/go-pg/pg/v10/orm"
)

type applyFn = func(q *orm.Query) (*orm.Query, error)
