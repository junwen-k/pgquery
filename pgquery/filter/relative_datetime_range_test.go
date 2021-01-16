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

var _ = Describe("RelativeDatetimeRange", func() {

	type RelativeDatetimeRangeTestItem struct {
		Id        int64
		Name      string
		CreatedAt time.Time `pg:"type:timestamp"` // Without timezone.
	}

	Context("marshalling json", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		When("layout is not set", func() {
			It("should marshal json with default format", func() {
				f := filter.NewRelativeDateTimeRange(time.Time{})
				f.At = &t

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`{"at":"2021-01-15T00:00:00Z","ago":{},"upcoming":{}}`))
			})
		})

		When("layout is set", func() {
			It("should marshal json with specified format", func() {
				f := filter.NewRelativeDateTimeRange(time.Time{}).Layout(time.Kitchen)
				f.At = &t

				b, err := json.Marshal(f)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`{"at":"12:00AM","ago":{},"upcoming":{}}`))
			})
		})
	})

	Context("unmarshalling json", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		When("layout is not set", func() {
			It("should unmarshal json with default format", func() {
				f := filter.NewRelativeDateTimeRange(time.Time{})
				f.At = &t

				err = json.Unmarshal([]byte(`{"at":"2021-01-15T00:00:00Z","ago":{},"upcoming":{}}`), f)
				Expect(err).ToNot(HaveOccurred())

				equal := filter.NewRelativeDateTimeRange(time.Time{}).Layout(time.RFC3339)
				equal.At = &t
				Expect(f).To(Equal(equal))
			})
		})

		When("layout is set", func() {
			It("should unmarshal json successfully", func() {
				f := filter.NewRelativeDateTimeRange(time.Time{}).Layout("2006-01-02")
				f.At = &t

				err := json.Unmarshal([]byte(`{"at":"2021-01-15"}`), f)
				Expect(err).ToNot(HaveOccurred())

				equal := filter.NewRelativeDateTimeRange(time.Time{}).Layout("2006-01-02")
				equal.At = &t
				Expect(f).To(Equal(equal))
			})
		})
	})

	Context("generating sql", func() {
		t, err := time.Parse("2006-01-02", "2021-01-15")
		Expect(err).ToNot(HaveOccurred())

		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &RelativeDatetimeRangeTestItem{})

			f := filter.NewRelativeDateTimeRange(time.Time{}).Column("created_at").AgoHour(5)
			f.At = &t
			q = f.Build(q.WhereGroup)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "relative_datetime_range_test_item"."id", "relative_datetime_range_test_item"."name", "relative_datetime_range_test_item"."created_at" FROM "relative_datetime_range_test_items" AS "relative_datetime_range_test_item" WHERE (("created_at" >= '2021-01-15T00:00:00Z'::timestamp - interval '5 hours') AND ("created_at" <= '2021-01-15T00:00:00Z'::timestamp))`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*RelativeDatetimeRangeTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &RelativeDatetimeRangeTestItem{
				Name: fmt.Sprintf("name-%d", itemCount),
			}
			if itemCount%2 == 0 {
				item.CreatedAt = testTime.Add(time.Duration(itemCount) * time.Hour)
			} else {
				item.CreatedAt = testTime.Add(-time.Duration(itemCount) * time.Hour)
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with hour ago filter", func() {
			var items []RelativeDatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewRelativeDateTimeRange(testTime).Column("created_at").AgoHour(5).Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(3)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					switch idx {
					case 0:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Second()))
					case 1:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Second()))
					case 2:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Second()))
					}
				}
			}
		})

		It("works with hour upcoming filter", func() {
			var items []RelativeDatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewRelativeDateTimeRange(testTime).Column("created_at").UpcomingHour(5).Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(2)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					switch idx {
					case 0:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Second()))
					case 1:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Second()))
					}
				}
			}
		})

		It("works with hour ago / upcoming filter", func() {
			var items []RelativeDatetimeRangeTestItem
			q := db.Model(&items)

			filter.NewRelativeDateTimeRange(testTime).Column("created_at").AgoHour(5).UpcomingHour(5).Build(q.WhereGroup)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(5)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					switch idx {
					case 0:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+2) * time.Hour).Second()))
					case 1:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+3) * time.Hour).Second()))
					case 2:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+4) * time.Hour).Second()))
					case 3:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(time.Duration(idx+5) * time.Hour).Second()))
					case 4:
						Expect(item.CreatedAt.Year()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Year()))
						Expect(item.CreatedAt.Month()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Month()))
						Expect(item.CreatedAt.Day()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Day()))
						Expect(item.CreatedAt.Second()).To(Equal(testTime.Add(-time.Duration(idx+6) * time.Hour).Second()))
					}
				}
			}
		})
	})
})
