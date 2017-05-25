package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
)

type Persistent interface {
	collection() string
}

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

func getId(id string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id)}
}

func findOne(document Persistent, id bson.M) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	return s.DB(database).C(document.collection()).Find(id).One(document)
}

func findQuery(collection string, result interface{}, query interface{}) error {
	s, err := getSession()
	if err != nil {
		return err
	}

	defer closeSession(s)

	return s.DB(database).C(collection).Find(query).All(result)
}

func findAll(collection string, result interface{}) error {
	return findQuery(collection, result, nil)
}

func insert(document Persistent) error {
	s, err := getSession()
	if err != nil {
		log.Println("fuck off")
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(document.collection()).Insert(document)
}

func remove(document Persistent, id bson.M) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(document.collection()).Remove(id)
}
