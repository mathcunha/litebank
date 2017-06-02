package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
)

const database = "litebank"

var (
	mongo_port = "localhost"
)

func init() {
	if m := os.Getenv("MONGO_PORT"); m != "" {
		mongo_port = strings.Replace(m, "tcp", "mongodb", 1)
		log.Printf("INFO: Mongo broker at %v \n", m)
	}
}

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial(mongo_port)
	return session, err
}

func closeSession(s *mgo.Session) {
	s.Close()
}

func findOne(collection string, result interface{}, id string) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)
	log.Printf("_id:%q", id)

	return s.DB(database).C(collection).Find(bson.M{"_id": id}).One(result)
}

func findQuery(col string, r interface{}, query interface{}) error {
	s, err := getSession()
	if err != nil {
		return err
	}

	defer closeSession(s)

	return s.DB(database).C(col).Find(query).All(r)
}

func findAll(col string, r interface{}) error {
	return findQuery(col, r, nil)
}

func insert(col string, doc interface{}) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(col).Insert(doc)
}

func remove(col string, id string) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(col).Remove(id)
}
