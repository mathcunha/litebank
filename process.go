package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
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

func (c *Costumer) create() (e *Event, err error) {
	c.Id = uuid.NewV4().String()
	c.Creation = time.Now()
	return newConsumerEvent(c)
}

func (a *Account) newEntity() Entity {
	return &Account{}
}

func (a *Account) create() (e *Event, err error) {
	a.Id = uuid.NewV4().String()
	return newAccountEvent(a)
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
func (e *Event) loadEntity() (Entity, error) {
	dec := json.NewDecoder(strings.NewReader(e.Body))
	if (e.Type & TCostumer) == TCostumer {
		c := Costumer{}
		if err := dec.Decode(&c); err != nil {
			return nil, err
		}
		return &c, nil

	} else if (e.Type & TAccount) == TAccount {
		a := Account{}
		if err := dec.Decode(&a); err != nil {
			return nil, err
		}
		return &a, nil
	}
	return nil, errors.New(fmt.Sprintf("Event type out of range %d", e.Type))

}
