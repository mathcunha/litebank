package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

func newAccount(c Costumer, number string) (a *Account) {
	a = &Account{Id: uuid.NewV4().String(), Number: number, Costumer: c}
	return
}

func newConsumer(name string) (c *Costumer) {
	c = &Costumer{Name: name, Id: uuid.NewV4().String(), Creation: time.Now()}
	return
}

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

func createAccount(c Costumer, number string) (a *Account, err error) {
	a = newAccount(c, number)
	newAccountEvent(a)
	return a, nil
}

func createConsumer(name string) (c *Costumer, err error) {
	c = newConsumer(name)
	newConsumerEvent(c)
	return c, nil
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
