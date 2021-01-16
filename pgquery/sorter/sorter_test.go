// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package sorter_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pgext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

func TestSorter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sorter Suite")
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
