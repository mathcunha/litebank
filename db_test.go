package main

import (
	"testing"
)

func TestFindOne(t *testing.T) {
	//func findOne(document Persistent, id string) error {
	c := Costumer{}
	id := "f67de2c9-8bf9-4093-b5fa-e2c230569374"
	if err := findOne(c.collection(), &c, id); err != nil {
		t.Error(err)
	}
	t.Logf("%v", c)
}
