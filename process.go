package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

func newEvent(v interface{}, eventType EventType) (e *Event, err error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	e = &Event{Type: eventType, Payload: body}
	fmt.Printf("Event(%v) - Payload (%v)\n", eventType, string(body))

	return e, nil
}

func (t *Transaction) newEntity() Entity {
	return &Transaction{}
}

func (t *Transaction) create() (e *Event, err error) {
	t.Id = uuid.NewV4().String()
	return newEvent(t, NewTransactionEvent)
}

func (c *Costumer) newEntity() Entity {
	return &Costumer{}
}

func (c *Costumer) create() (e *Event, err error) {
	c.Id = uuid.NewV4().String()
	c.Creation = time.Now()
	return newEvent(c, NewCostumerEvent)
}

func (a *Account) newEntity() Entity {
	return &Account{}
}

func (a *Account) create() (e *Event, err error) {
	a.Id = uuid.NewV4().String()
	return newEvent(a, NewAccountEvent)
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
	} else if (e.Type & TTransaction) == TTransaction {
		t := Transaction{}
		if err := json.Unmarshal(e.Payload, &t); err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, errors.New(fmt.Sprintf("Event type out of range %d", e.Type))

}
