// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package sorter

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// OrderDirection order direction enum type.
type OrderDirection int

const (
	// OrderDirectionAsc order direction ascending enum.
	OrderDirectionAsc OrderDirection = iota

	// OrderDirectionDesc order direction descending enum.
	OrderDirectionDesc
)

// String returns the string presentation for the order direction.
func (d OrderDirection) String() string {
	return [...]string{"ASC", "DESC"}[d]
}

// Order order common sorter.
type Order struct {
	column    string
	Direction OrderDirection `json:"direction,omitempty"`
}

// MarshalJSON custom JSON marshaler.
func (s *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Direction.String())
}

// UnmarshalJSON custom JSON unmarshaler.
func (s *Order) UnmarshalJSON(b []byte) error {
	type alias Order

	m1 := struct {
		Direction string `json:"direction,omitempty"`
		*alias
	}{alias: (*alias)(s)}
	var m2 string

	if err := json.Unmarshal(b, &m1); err == nil {
		d := strings.ToLower(m1.Direction)
		if d == strings.ToLower(OrderDirectionAsc.String()) {
			s.Direction = OrderDirectionAsc
		}
		if d == strings.ToLower(OrderDirectionDesc.String()) {
			s.Direction = OrderDirectionDesc
		}
		return nil
	}

	if err := json.Unmarshal(b, &m2); err == nil {
		d := strings.ToLower(m2)
		if d == strings.ToLower(OrderDirectionAsc.String()) {
			s.Direction = OrderDirectionAsc
		}
		if d == strings.ToLower(OrderDirectionDesc.String()) {
			s.Direction = OrderDirectionDesc
		}
		return nil
	}

	return nil // TODO: return unsupported format error
}

// NewOrder initializes a new order sorter.
func NewOrder(direction OrderDirection) *Order {
	return &Order{
		Direction: direction,
	}
}

// NewOrderAsc initializes a new ascending order sorter.
func NewOrderAsc() *Order {
	o := &Order{}
	return o.Asc()
}

// NewOrderDesc initializes a new order sorter.
func NewOrderDesc() *Order {
	o := &Order{}
	return o.Desc()
}

// Column sets the column for the order sorter.
func (s *Order) Column(column string) *Order {
	s.column = column
	return s
}

// Asc sets the direction to ascending order.
func (s *Order) Asc() *Order {
	s.Direction = OrderDirectionAsc
	return s
}

// Desc sets the direction to descending order.
func (s *Order) Desc() *Order {
	s.Direction = OrderDirectionDesc
	return s
}

// Build build query.
func (s *Order) Build(q *orm.Query) *orm.Query {
	return q.OrderExpr("? ?", types.Ident(s.column), types.Safe(s.Direction.String()))
}
