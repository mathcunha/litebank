package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

func newAccountEvent(a *Account) (e *Event, err error) {
	return newEvent(*a, NewAccountEvent)
}

func newConsumerEvent(c *Costumer) (e *Event, err error) {
	return newEvent(*c, NewCostumerEvent)
}

func newEvent(v interface{}, eventType EventType) (e *Event, err error) {
	body, err := getJson(v)
	if err != nil {
		return nil, err
	}

	e = &Event{Type: eventType, Body: body}
	fmt.Printf("Event(%v) - body (%v)\n", eventType, body)

	return e, nil
}

func (c *Costumer) newEntity() Entity {
	return &Costumer{}
}

func (c *Costumer) create() (err error) {
	c.Id = uuid.NewV4().String()
	c.Creation = time.Now()
	_, err = newConsumerEvent(c)
	return err
}

func (a *Account) newEntity() Entity {
	return &Account{}
}

func (a *Account) create() (err error) {
	a.Id = uuid.NewV4().String()
	_, err = newAccountEvent(a)
	return err
}

func getJson(v interface{}) (body string, e error) {
	buffer := bytes.NewBufferString(body)
	enc := json.NewEncoder(buffer)
	e = enc.Encode(v)
	if e != nil {
		return body, e
	}
	return buffer.String(), nil
}
