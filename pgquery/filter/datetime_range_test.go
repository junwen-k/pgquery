// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package filter_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/filter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatetimeRange", func() {

	type DatetimeRangeTestItem struct {
		Id        int64
		Name      string
		CreatedAt time.Time
	}

	Context("marshalling json", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		When("layout is not set", func() {
			It("should marshal json with default format", func() {
				f := filter.NewDateTimeRange().After(t)

				b, err := json.Marshal(f)
				Expect(err).ToNot(HaveOccurred())

				Expect(b).To(MatchJSON(`{"after":"2021-01-15T00:00:00Z"}`))
			})
		})

		When("layout is set", func() {
			It("should marshal json with specified format", func() {
				f := filter.NewDateTimeRange().Layout(time.Kitchen).After(t)

				b, err := json.Marshal(f)
				Expect(err).ToNot(HaveOccurred())

				Expect(b).To(MatchJSON(`{"after":"12:00AM"}`))
			})
		})
	})

	Context("unmarshalling json", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		When("layout is not set", func() {
			It("should unmarshal json with default format", func() {
				f := filter.NewDateTimeRange()

				err = json.Unmarshal([]byte(`{"after":"2021-01-15T00:00:00Z"}`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(filter.NewDateTimeRange().After(t).Layout(time.RFC3339)))
			})
		})

		When("layout is set", func() {
			It("should unmarshal json with specified format", func() {
				f := filter.NewDateTimeRange().Layout("2006-01-02")

				err = json.Unmarshal([]byte(`{"after":"2021-01-15"}`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(filter.NewDateTimeRange().After(t).Layout("2006-01-02")))
			})
		})

	})

	Context("generating sql", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &DatetimeRangeTestItem{})

			q = filter.NewDateTimeRange().Column("created_at").After(t).Build(q.WhereGroup)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "datetime_range_test_item"."id", "datetime_range_test_item"."name", "datetime_range_test_item"."created_at" FROM "datetime_range_test_items" AS "datetime_range_test_item" WHERE (("created_at" > '2021-01-15T00:00:00Z'))`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*DatetimeRangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &DatetimeRangeTestItem{
				Name:      fmt.Sprintf("name-%d", itemCount),
				CreatedAt: testTime.Add(time.Duration(itemCount) * time.Hour),
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with after filter", func() {
			var items []DatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewDateTimeRange().After(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+6)))
					Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+6) * time.Hour).Year()))
					Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+6) * time.Hour).Month()))
					Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+6) * time.Hour).Day()))
					Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+6) * time.Hour).Second()))
				}
			}
		})

		It("works with before filter", func() {
			var items []DatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewDateTimeRange().Before(testTime.Add(6 * time.Hour)).Column("created_at").Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
					Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Year()))
					Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Month()))
					Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Day()))
					Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Second()))
				}
			}
		})

		It("works with from filter", func() {
			var items []DatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewDateTimeRange().From(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(6)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+5)))
					Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Year()))
					Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Month()))
					Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Day()))
					Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Second()))
				}
			}
		})

		It("works with to filter", func() {
			var items []DatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewDateTimeRange().To(testTime.Add(5 * time.Hour)).Column("created_at").Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
					Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Year()))
					Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Month()))
					Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Day()))
					Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+1) * time.Hour).Second()))
				}
			}
		})

		It("works with complex filter", func() {
			var items []DatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewDateTimeRange().From(testTime.Add(2 * time.Hour)).To(testTime.Add(6 * time.Hour)).Column("created_at").Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+2)))
					Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+2) * time.Hour).Year()))
					Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+2) * time.Hour).Month()))
					Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+2) * time.Hour).Day()))
					Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+2) * time.Hour).Second()))
				}
			}
		})
	})
})
