// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pgext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}

func selectQueryString(q *orm.Query) string {
	sel := orm.NewSelectQuery(q)
	s := queryString(sel)
	return s
}

func queryString(f orm.QueryAppender) string {
	fmter := orm.NewFormatter().WithModel(f)
	b, err := f.AppendQuery(fmter, nil)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
