package storage

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
