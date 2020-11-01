// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"os"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pgext"
)

var (
	db       *pg.DB
	testTime time.Time
)

func init() {
	var err error
	db = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "password",
	})
	if ok, _ := strconv.ParseBool(os.Getenv("DB_ENABLE_LOGGING")); ok {
		db.AddQueryHook(pgext.DebugHook{Verbose: true})
	}
	testTime, err = time.Parse(time.RFC822, "31 Oct 20 00:00 UTC")
	if err != nil {
		panic(err)
	}
}
