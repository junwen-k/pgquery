// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/filter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Match", func() {

	type MatchTestItem struct {
		Id   int64
		Name string
	}

	Context("marshalling json", func() {
		When("using a single value", func() {
			It("should marshal json successfully", func() {
				f := filter.NewMatch("match")

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`"match"`))
			})
		})

		When("using an array of values", func() {
			It("should marshal json successfully", func() {
				f := filter.NewMatch("match_1", "match_2")

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`["match_1","match_2"]`))
			})
		})
	})

	Context("unmarshalling json", func() {
		When("using object syntax", func() {
			When("using a single value", func() {
				It("should unmarshal json successfully", func() {
					f := filter.NewMatch()

					err := json.Unmarshal([]byte(`{"values":"match"}`), f)
					Expect(err).ToNot(HaveOccurred())

					Expect(f).To(Equal(filter.NewMatch("match")))
				})
			})

			When("using an array of values", func() {
				It("should unmarshal json successfully", func() {
					f := filter.NewMatch()

					err := json.Unmarshal([]byte(`{"values":["match_1","match_2"]}`), f)
					Expect(err).ToNot(HaveOccurred())

					Expect(f).To(Equal(filter.NewMatch("match_1", "match_2")))
				})
			})
		})

		When("using non-object syntax", func() {
			When("using an array of values", func() {
				It("should unmarshal json successfully", func() {
					f := filter.NewMatch()

					err := json.Unmarshal([]byte(`["match_1","match_2"]`), f)
					Expect(err).ToNot(HaveOccurred())

					Expect(f).To(Equal(filter.NewMatch("match_1", "match_2")))
				})
			})
		})

		When("using a single value", func() {
			It("should unmarshal json successfully", func() {
				f := filter.NewMatch()

				err := json.Unmarshal([]byte(`"match"`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(filter.NewMatch("match")))
			})
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &MatchTestItem{})

			q = filter.NewMatch("match").Column("name").Build(q.Where)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "match_test_item"."id", "match_test_item"."name" FROM "match_test_items" AS "match_test_item" WHERE ("name" = 'match')`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*MatchTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &MatchTestItem{
				Name: fmt.Sprintf("name-%d", itemCount),
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with single value", func() {
			var items []MatchTestItem
			q := db.Model(&items)

			filter.NewMatch("name-1").Column("name").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(1)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
			}
		})

		It("works with multiple values", func() {
			var items []MatchTestItem
			q := db.Model(&items)

			filter.NewMatch("name-1", "name-2", "name-3").Column("name").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(3)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
				}
			}
		})
	})
})
