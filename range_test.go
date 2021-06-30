// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package pgquery_test

import (
	"encoding/json"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/pgquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Range", func() {

	type RangeTestItem struct {
		Id     int64
		Age    int
		Height int
	}

	Context("marshalling json", func() {
		It("should marshal json successfully", func() {
			f := pgquery.NewRange("").GreaterThan(0)

			b, err := json.Marshal(f)
			Expect(err).NotTo(HaveOccurred())

			Expect(b).To(MatchJSON(`{"gt":0}`))
		})
	})

	Context("unmarshalling json", func() {
		It("should unmarshal json successfully", func() {
			f := pgquery.NewRange("")

			err := json.Unmarshal([]byte(`{"gt":0}`), f)
			Expect(err).ToNot(HaveOccurred())

			Expect(f).To(Equal(pgquery.NewRange("").GreaterThan(0)))
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &RangeTestItem{})

			q.WhereGroup(pgquery.NewRange("age").GreaterThan(0).Appender())

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "range_test_item"."id", "range_test_item"."age", "range_test_item"."height" FROM "range_test_items" AS "range_test_item" WHERE (("age" > 0))`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*RangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &RangeTestItem{
				Age:    itemCount,
				Height: (itemCount * 10) + 130,
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with greater than filter", func() {
			var items []RangeTestItem
			q := db.Model(&items)

			q.WhereGroup(pgquery.NewRange("age").GreaterThan(5).Appender())

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Age).To(Equal(idx + 6))
				}
			}
		})

		It("works with greater than equal filter", func() {
			var items []RangeTestItem
			q := db.Model(&items)

			q.WhereGroup(pgquery.NewRange("age").GreaterThanEqual(5).Appender())

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(6)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Age).To(Equal(idx + 5))
				}
			}
		})

		It("works with less than filter", func() {
			var items []RangeTestItem
			q := db.Model(&items)

			q.WhereGroup(pgquery.NewRange("age").LessThan(5).Appender())

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(4)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Age).To(Equal(idx + 1))
				}
			}
		})

		It("works with less than equal filter", func() {
			var items []RangeTestItem
			q := db.Model(&items)

			q.WhereGroup(pgquery.NewRange("age").LessThanEqual(5).Appender())

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Age).To(Equal(idx + 1))
				}
			}
		})

		It("works with complex filter", func() {
			var items []RangeTestItem
			q := db.Model(&items)

			q.WhereGroup(pgquery.NewRange("age").GreaterThan(5).LessThan(8).Appender())
			q.WhereOrGroup(pgquery.NewRange("height").GreaterThan(5).LessThan(8).Appender())

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(2)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Age).To(Equal(idx + 6))
				}
			}
		})
	})
})
