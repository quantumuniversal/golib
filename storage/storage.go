package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Query --
type Query struct {
	Column     string
	ColumnExpr string // count(*) as column_count
	Where      []Where
	WhereOr    []Where
	WhereIn    []Where
	Relations  []string
	Order      string
	Limit      int
}

// MgoQuery --
type MgoQuery struct {
	Index mgo.Index
	One   bool
	Where bson.M
	Sort  string
	Limit int
}
