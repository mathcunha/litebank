package main

import (
	"encoding/json"
	"errors"
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

	e = &Event{Type: eventType, Payload: body}
	fmt.Printf("Event(%v) - Payload (%v)\n", eventType, string(body))

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

func getJson(v interface{}) ([]byte, error) {
	body, e := json.Marshal(v)
	if e != nil {
		return body, e
	}
	return body, nil
}
func (e *Event) loadEntity() (Entity, error) {
	if (e.Type & TCostumer) == TCostumer {
		c := Costumer{}
		if err := json.Unmarshal(e.Payload, &c); err != nil {
			return nil, err
		}
		return &c, nil

	} else if (e.Type & TAccount) == TAccount {
		a := Account{}
		if err := json.Unmarshal(e.Payload, &a); err != nil {
			return nil, err
		}
		return &a, nil
	}
	return nil, errors.New(fmt.Sprintf("Event type out of range %d", e.Type))

}
