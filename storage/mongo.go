package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/quantumuniversal/golib/sec"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo --
type Mongo struct {
	Auth       MgoAuth
	DB         string
	Collection string
	Query      MgoQuery
	Data       interface{}
}

// MgoAuth --
type MgoAuth struct {
	Host     string
	Username string
	Password string
	Port     string
}

func (m *Mongo) executeQuery(c *mgo.Collection) error {
	var err error
	if m.Query.One {
		err = c.Find(m.Query.Where).Sort(m.Query.Sort).One(&m.Data)

		if err != nil {
			return err
		}

	} else {
		if m.Query.Limit > 0 {
			err = c.Find(m.Query.Where).Sort(m.Query.Sort).Limit(m.Query.Limit).All(&m.Data)
		} else {
			err = c.Find(m.Query.Where).Sort(m.Query.Sort).All(&m.Data)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Session --
func (m *Mongo) Session() (*mgo.Session, *mgo.Collection, error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s@%s:%s", m.Auth.Username, m.Auth.Password, m.Auth.Host, m.Auth.Port))

	if err != nil {
		return &mgo.Session{}, &mgo.Collection{}, err
	}

	collection := session.DB(fmt.Sprintf("%s", m.DB)).C(m.Collection)

	return session, collection, nil
}

// CreateUser --
func (m *Mongo) CreateUser(accountNumber string, services []string, authDB string) (map[string]interface{}, error) {
	userData := make(map[string]interface{})

	session, err := mgo.Dial(fmt.Sprintf("%s:%s@%s:%s", m.Auth.Username, m.Auth.Password, m.Auth.Host, m.Auth.Port))

	if err != nil {
		return nil, err
	}

	var db mgo.Database
	db.Session = session
	db.Name = authDB

	var user mgo.User
	h := sha1.New()
	h.Write([]byte(sec.RandomString(16)))

	user.Username = hex.EncodeToString(h.Sum(nil))

	h = sha1.New()
	h.Write([]byte(sec.RandomString(16)))
	user.Password = hex.EncodeToString(h.Sum(nil))

	var role []mgo.Role
	role = append(role, "readWrite")

	user.OtherDBRoles = make(map[string][]mgo.Role)
	for idx := range services {
		serviceName := fmt.Sprintf("%s-%s", accountNumber, services[idx])
		user.OtherDBRoles[serviceName] = role
	}

	err = db.UpsertUser(&user)

	if err != nil {
		return nil, err
	}

	userData["username"] = user.Username
	userData["password"] = user.Password

	if err != nil {
		return nil, err
	}

	return userData, nil
}

// Select --
func (m *Mongo) Select() error {
	session, collection, err := m.Session()

	if err != nil {
		return err
	}

	err = m.executeQuery(collection)

	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}

// Insert --
func (m *Mongo) Insert() error {
	session, collection, err := m.Session()

	if err != nil {
		return err
	}

	err = collection.EnsureIndex(m.Query.Index)

	if err != nil {
		return err
	}

	err = collection.Insert(&m.Data)

	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}

// Update --
func (m *Mongo) Update(condition bson.M) error {
	session, collection, err := m.Session()

	if err != nil {
		return err
	}

	mappedData := m.Data.(map[string]interface{})

	err = collection.Update(condition, bson.M{"$set": bson.M(mappedData)})

	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}

// Delete --
func (m *Mongo) Delete(condition bson.M) error {
	session, collection, err := m.Session()

	if err != nil {
		return err
	}

	err = collection.Remove(condition)

	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}
