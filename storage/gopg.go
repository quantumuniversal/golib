package storage

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
	"github.com/jinzhu/inflection"
)

// Gopg --
type Gopg struct {
	Auth      pg.Options
	Schema    string
	TableName string
	Query     Query
	Model     interface{}
	Data      interface{}
}

// Query --
type Query struct {
	Column     string
	ColumnExpr string // count(*) as column_count
	Where      []Where
	WhereOr    []Where
	Relation   string
	Order      string
	Limit      int
}

// Where --
type Where struct {
	Condition string
	Param     interface{}
}

func (g *Gopg) prepareStatement(q *orm.Query) *orm.Query {
	if g.Query.Column != "" {
		q = q.Column(g.Query.Column)
	}

	if g.Query.ColumnExpr != "" {
		q = q.ColumnExpr(g.Query.ColumnExpr)
	}

	if g.Query.Relation != "" {
		q = q.Relation(g.Query.Relation)
	}

	if len(g.Query.Where) > 0 {
		for _, statement := range g.Query.Where {
			q = q.Where(statement.Condition, statement.Param)
		}
	}

	if len(g.Query.WhereOr) > 0 {
		for _, statement := range g.Query.WhereOr {
			q = q.WhereOr(statement.Condition, statement.Param)
		}
	}

	if g.Query.Order != "" {
		q = q.Order(g.Query.Order)
	}

	if g.Query.Limit != 0 {
		q = q.Limit(g.Query.Limit)
	}

	return q
}

// Init --
func (g *Gopg) Init(models []interface{}) error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	_, err := db.Exec(fmt.Sprintf("CREATE Schema IF NOT EXISTS \"%s\";", g.Schema))

	if err != nil {
		return err
	}

	for _, model := range models {
		orm.SetTableNameInflector(func(s string) string {
			return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
		})
		db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
	}

	return nil
}

// Select --
func (g *Gopg) Select() error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
	})

	q := g.prepareStatement(db.Model(g.Data))

	err := q.Select()

	if err != nil {
		return err
	}

	return nil
}

// SelectCount --
func (g *Gopg) SelectCount() (int, error) {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
	})

	q := g.prepareStatement(db.Model(g.Data))

	count, err := q.Count()

	if err != nil {
		return 0, err
	}

	return count, nil
}

// Insert --
func (g *Gopg) Insert() error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
	})

	err := db.Insert(g.Data)

	if err != nil {
		return err
	}

	return nil
}

// Update --
func (g *Gopg) Update(uid uuid.UUID) error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
	})

	_, err := db.Model(g.Data).Where("id = ?", uid).Update()

	if err != nil {
		return err
	}

	return nil
}

// Delete --
func (g *Gopg) Delete(condition string, param interface{}) error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(s))
	})

	_, err := db.Model(g.Data).Where(condition, param).Delete()
	if err != nil {
		return err
	}

	return nil
}
