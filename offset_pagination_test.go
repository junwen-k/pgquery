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

var _ = Describe("OffsetPagination", func() {

	type OffsetPaginationTestItem struct {
		Id   int64
		Name string
	}

	Context("marshalling json", func() {
		It("should marshal json successfully", func() {
			f := pgquery.NewOffsetPagination().Offset(1, 10)

			b, err := json.Marshal(f)
			Expect(err).NotTo(HaveOccurred())

			Expect(b).To(MatchJSON(`{"page":1,"limit":10}`))
		})
	})

	Context("unmarshalling json", func() {
		It("should unmarshal json successfully", func() {
			f := pgquery.NewOffsetPagination().Offset(0, 0)

			err := json.Unmarshal([]byte(`{"page":1,"limit":10}`), f)
			Expect(err).ToNot(HaveOccurred())

			Expect(f).To(Equal(pgquery.NewOffsetPagination().Offset(1, 10)))
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &OffsetPaginationTestItem{})

			q = pgquery.NewOffsetPagination().Offset(1, 10).Build(q)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "offset_pagination_test_item"."id", "offset_pagination_test_item"."name" FROM "offset_pagination_test_items" AS "offset_pagination_test_item" LIMIT 10`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*OffsetPaginationTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &OffsetPaginationTestItem{
				Name: fmt.Sprintf("name-%d", itemCount),
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with page and limit", func() {
			var items []OffsetPaginationTestItem
			q := db.Model(&items)

			pgquery.NewOffsetPagination().Offset(1, 5).Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
				}
			}
		})

		It("works with no page", func() {
			var items []OffsetPaginationTestItem
			q := db.Model(&items)

			pgquery.NewOffsetPagination().Offset(0, 5).Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
				}
			}
		})

		It("works with no limit", func() {
			var items []OffsetPaginationTestItem
			q := db.Model(&items)

			pgquery.NewOffsetPagination().Offset(1, 0).Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(10)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
				}
			}
		})
	})
})
