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

var _ = Describe("KeywordSearch", func() {

	type KeywordSearchTestItem struct {
		Id     int64
		Name   string
		Emails []string `pg:",array"`
	}

	Context("marshalling json", func() {
		It("should marshal json successfully", func() {
			f := pgquery.NewKeywordSearch("").Keyword("keyword")

			b, err := json.Marshal(f)
			Expect(err).NotTo(HaveOccurred())

			Expect(b).To(MatchJSON(`"keyword"`))
		})
	})

	Context("unmarshalling json", func() {
		When("using object syntax", func() {
			It("should unmarshal json successfully", func() {
				f := pgquery.NewKeywordSearch("")

				err := json.Unmarshal([]byte(`{"value":"keyword"}`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(pgquery.NewKeywordSearch("").Keyword("keyword")))
			})
		})

		When("using non-object syntax", func() {
			It("should unmarshal json successfully", func() {
				f := pgquery.NewKeywordSearch("")

				err := json.Unmarshal([]byte(`"keyword"`), f)
				Expect(err).ToNot(HaveOccurred())

				Expect(f).To(Equal(pgquery.NewKeywordSearch("").Keyword("keyword")))
			})
		})
	})

	Context("generating sql", func() {
		It("should generate correct SQL string", func() {
			q := orm.NewQuery(nil, &KeywordSearchTestItem{})

			q = pgquery.NewKeywordSearch("name").Keyword("keyword").Build(q.Where)

			s := queryString(q)
			Expect(s).To(Equal(`SELECT "keyword_search_test_item"."id", "keyword_search_test_item"."name", "keyword_search_test_item"."emails" FROM "keyword_search_test_items" AS "keyword_search_test_item" WHERE ("name" LIKE '%keyword%')`))
		})
	})

	Context("integration testing", func() {
		err := db.Model((*KeywordSearchTestItem)(nil)).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		Expect(err).ToNot(HaveOccurred())

		for itemCount := 1; itemCount <= 10; itemCount++ {
			item := &KeywordSearchTestItem{
				Name:   fmt.Sprintf("name-%d", itemCount),
				Emails: make([]string, 0),
			}
			for emailCount := 1; emailCount <= 5; emailCount++ {
				item.Emails = append(item.Emails, fmt.Sprintf("email-%d(%d)@root", itemCount, emailCount))
			}
			_, err = db.Model(item).Insert()
			Expect(err).ToNot(HaveOccurred())
		}

		It("works with default search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").Keyword("name-1").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(2)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
				Expect(items[0].Emails).To(Equal([]string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}))

				Expect(items[1].Id).ToNot(BeZero())
				Expect(items[1].Name).To(Equal("name-10"))
				Expect(items[1].Emails).To(Equal([]string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}))
			}
		})

		It("works with multiple columns search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").Keyword("(1)@root").Build(q.Where)
			pgquery.NewKeywordSearch("emails,array").Keyword("(1)@root").Build(q.WhereOr)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(10)) {
				for idx, item := range items {
					Expect(item.Id).ToNot(BeZero())
					Expect(item.Name).To(Equal(fmt.Sprintf("name-%d", idx+1)))
					emails := make([]string, 0)
					for emailCount := 1; emailCount <= 5; emailCount++ {
						emails = append(emails, fmt.Sprintf("email-%d(%d)@root", idx+1, emailCount))
					}
					Expect(item.Emails).To(Equal(emails))
				}
			}
		})

		It("works with match all search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").MatchAll().Keyword("name-1").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(1)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
				Expect(items[0].Emails).To(Equal([]string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}))
			}
		})

		It("works with match start search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").MatchStart().Keyword("name-1").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(2)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
				Expect(items[0].Emails).To(Equal([]string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}))

				Expect(items[1].Id).ToNot(BeZero())
				Expect(items[1].Name).To(Equal("name-10"))
				Expect(items[1].Emails).To(Equal([]string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}))
			}
		})

		It("works with match end search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").MatchEnd().Keyword("-1").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(1)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
				Expect(items[0].Emails).To(Equal([]string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}))
			}
		})

		It("works with case insensitive search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("name").CaseInsensitive().Keyword("NAME-10").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(1)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-10"))
				Expect(items[0].Emails).To(Equal([]string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}))
			}
		})

		It("works with array search", func() {
			var items []KeywordSearchTestItem
			q := db.Model(&items)

			pgquery.NewKeywordSearch("emails,array").Keyword("email-1").Build(q.Where)

			err := q.Select()
			Expect(err).ToNot(HaveOccurred())

			if Expect(items).To(HaveLen(2)) {
				Expect(items[0].Id).ToNot(BeZero())
				Expect(items[0].Name).To(Equal("name-1"))
				Expect(items[0].Emails).To(Equal([]string{"email-1(1)@root", "email-1(2)@root", "email-1(3)@root", "email-1(4)@root", "email-1(5)@root"}))

				Expect(items[1].Id).ToNot(BeZero())
				Expect(items[1].Name).To(Equal("name-10"))
				Expect(items[1].Emails).To(Equal([]string{"email-10(1)@root", "email-10(2)@root", "email-10(3)@root", "email-10(4)@root", "email-10(5)@root"}))
			}
		})
	})
})
