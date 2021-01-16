// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package sorter_test

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/go-pg/pg/v10/orm"
	"github.com/junwen-k/go-pgquery/pgquery/sorter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order", func() {

	type OrderTestItem struct {
		Id   int64
		Name string
		Age  int
	}

	Context("marshalling json", func() {
		It("should marshal json successfully", func() {
			s := sorter.NewOrderAsc()

			b, err := json.Marshal(s)
			Expect(err).NotTo(HaveOccurred())

			Expect(b).To(MatchJSON(`"ASC"`))
		})
	})

	Context("unmarshalling json", func() {
		When("using object syntax", func() {
			When("using uppercase", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`{"direction":"DESC"}`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})

			When("using mixed case", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`{"direction":"DeSc"}`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})

			When("using lowercase", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`{"direction":"desc"}`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})
		})

		When("using non-object syntax", func() {
			When("using uppercase", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`"DESC"`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})

			When("using mixed case", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`"DeSc"`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})

			When("using lowercase", func() {
				It("should unmarshal json successfully", func() {
					s := sorter.NewOrderAsc()

					err := json.Unmarshal([]byte(`"desc"`), s)
					Expect(err).ToNot(HaveOccurred())

					Expect(s).To(Equal(sorter.NewOrderDesc()))
				})
			})
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &OrderTestItem{})

			q = sorter.NewOrderAsc().Column("age").Build(q)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "order_test_item"."id", "order_test_item"."name", "order_test_item"."age" FROM "order_test_items" AS "order_test_item" ORDER BY "age" ASC`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*OrderTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &OrderTestItem{
				Name: fmt.Sprintf("name-%d", itemCount),
				Age:  itemCount,
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with asc value", func() {
			var items []OrderTestItem
			q := db.Model(&items)

			sorter.NewOrderAsc().Column("age").Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(10)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
					Expect(item.Age).To(Equal(idx + 1))
				}
			}
		})

		It("works with desc value", func() {
			var items []OrderTestItem
			q := db.Model(&items)

			sorter.NewOrderDesc().Column("age").Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(10)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", len(items)-idx)))
					Expect(item.Age).To(Equal(len(items) - idx))
				}
			}
		})

		It("works with multiple columns asc value", func() {
			var items []OrderTestItem
			q := db.Model(&items)

			sorter.NewOrderDesc().Column("name").Build(q)
			sorter.NewOrderDesc().Column("age").Build(q)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(10)) {
				for idx, item := range items {
					switch idx {
					case 8:
						Expect(item.Id).ToNot(BeZero())
						Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", 10)))
						Expect(item.Age).To(Equal(10))
					case 9:
						Expect(item.Id).ToNot(BeZero())
						Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", 1)))
						Expect(item.Age).To(Equal(1))
					default:
						Expect(item.Id).ToNot(BeZero())
						Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", int(math.Abs(float64(idx-9))))))
						Expect(item.Age).To(Equal(int(math.Abs(float64(idx - 9)))))
					}
				}
			}
		})
	})
})
