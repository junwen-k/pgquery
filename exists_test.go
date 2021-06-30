// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/pgquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exists", func() {

	type ExistsTestItem struct {
		Id   int64
		Name string
	}

	Context("marshalling json", func() {
		When("value is nil", func() {
			It("should marshal json successfully", func() {
				f := pgquery.NewExists("name")

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`null`))
			})
		})

		When("value is false", func() {
			It("should marshal json successfully", func() {
				f := pgquery.NewExists("name").ShouldNotExists()

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`false`))
			})
		})
	})

	Context("unmarshalling json", func() {
		When("using object syntax", func() {
			It("should unmarshal json successfully", func() {
				f := pgquery.NewExists("name")

				err := json.Unmarshal([]byte(`{"value":true}`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(pgquery.NewExists("name").ShouldExists()))
			})
		})

		When("using non-object syntax", func() {
			It("should unmarshal json successfully", func() {
				f := pgquery.NewExists("name")

				err := json.Unmarshal([]byte(`true`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(pgquery.NewExists("name").ShouldExists()))
			})
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &ExistsTestItem{})

			q = pgquery.NewExists("name").ShouldExists().Build(q.Where)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "exists_test_item"."id", "exists_test_item"."name" FROM "exists_test_items" AS "exists_test_item" WHERE ("name" IS NOT NULL)`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*ExistsTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			name := ""
			if itemCount%2 == 0 {
				name = fmt.Sprintf("name-%d", itemCount)
			}
			item := &ExistsTestItem{
				Name: name,
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with true value", func() {
			var items []ExistsTestItem
			q := db.Model(&items)

			pgquery.NewExists("name").ShouldExists().Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", (idx+1)*2)))
				}
			}
		})

		It("works with false value", func() {
			var items []ExistsTestItem
			q := db.Model(&items)

			pgquery.NewExists("name").ShouldNotExists().Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for _, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(BeEmpty())
				}
			}
		})
	})
})
