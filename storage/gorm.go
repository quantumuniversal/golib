package storage

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
	"github.com/quantumuniversal/golib/enum/dbtype"
)

// Gorm --
type Gorm struct {
	Auth      Auth
	Schema    string
	TableName string
	Query     Query
	Model     interface{}
	Data      interface{}
}

// Auth --
type Auth struct {
	Type     string
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

// connect --
func (g *Gorm) connect() (*gorm.DB, error) {
	switch g.Auth.Type {
	case string(dbtype.POSTGRES):
		return gorm.Open(string(dbtype.POSTGRES),
			fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", g.Auth.Host, g.Auth.Port, g.Auth.User, g.Auth.DbName, g.Auth.Password))
	default:
		return nil, errors.New("Invalid Database Type")
	}

}

// setTableName --
func (g *Gorm) setTableName() string {
	return fmt.Sprintf("%s.%s", g.Schema, inflection.Plural(g.TableName))
}

// Init --
func (g *Gorm) Init(models []interface{}) error {
	db, err := g.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	db.AutoMigrate(models)

	return nil
}

// Select --
func (g *Gorm) Select() error {
	db, err := g.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	tableName := g.setTableName()

	db.Table(tableName).Find(g.Data)

	return nil
}

// Insert --
func (g *Gorm) Insert() error {
	db, err := g.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	tableName := g.setTableName()

	db.Table(tableName).Create(g.Data)

	return nil
}

// Update --
func (g *Gorm) Update(uid uuid.UUID) error {
	db, err := g.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	tableName := g.setTableName()

	db.Table(tableName).Model(g.Model).Updates(g.Data)

	return nil
}

// Delete --
func (g *Gorm) Delete(uid uuid.UUID) error {
	db, err := g.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	tableName := g.setTableName()

	db.Table(tableName).Where("id = ?", uid).Delete(g.Data)

	return nil
}
