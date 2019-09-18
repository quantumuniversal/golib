package storage

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
)

// Gopg --
type Gopg struct {
	Auth      pg.Options
	Schema    string
	TableName string
	Model     interface{}
	Data      interface{}
}

// SelectMany --
func (g *Gopg) SelectMany() error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, g.TableName)
	})

	err := db.Model(g.Data).Select()
	if err != nil {
		return err
	}

	return nil
}

// SelectOne --
func (g *Gopg) SelectOne(condition string, param interface{}) error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, g.TableName)
	})

	err := db.Model(g.Data).Where(condition, param).Select()
	if err != nil {
		return err
	}

	return nil
}

// Insert --
func (g *Gopg) Insert() error {
	db := pg.Connect(&g.Auth)
	defer db.Close()

	_, err := db.Exec(fmt.Sprintf("CREATE Schema IF NOT EXISTS \"%s\";", g.Schema))

	if err != nil {
		return err
	}

	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s.%s", g.Schema, g.TableName)
	})

	db.CreateTable(g.Model, &orm.CreateTableOptions{
		Temp: false,
	})

	err = db.Insert(g.Data)

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
		return fmt.Sprintf("%s.%s", g.Schema, g.TableName)
	})

	_, err := db.Model(g.Data).Where("id = ?", uid).Update()

	if err != nil {
		return err
	}

	return nil
}
