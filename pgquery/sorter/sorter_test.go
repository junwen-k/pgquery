// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package sorter_test

import (
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pgext"
)

var db *pg.DB

func init() {
	db = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
	})
	if ok, _ := strconv.ParseBool(os.Getenv("DB_ENABLE_LOGGING")); ok {
		db.AddQueryHook(pgext.DebugHook{Verbose: true})
	}
}
