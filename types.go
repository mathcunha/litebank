package main

import (
	"time"
)

type Entity interface {
	create() (*Event, error)
	newEntity() Entity
	collection() string
}

type Costumer struct {
	Id           string        `bson:"_id"`
	Name         string        `bson:",omitempty"`
	Creation     time.Time     `bson:",omitempty"`
	Accounts     []Account     `bson:",omitempty"`
	Transactions []Transaction `bson:",omitempty"`
}

func (c *Costumer) collection() string {
	return "costumer"
}

type Account struct {
	Id       string `bson:"_id"`
	Number   string
	Balance  float64
	Costumer Costumer
	Queue    []float64
}

func (a *Account) collection() string {
	return "account"
}

func (c *Transaction) collection() string {
	return "transaction"
}

type Event struct {
	Type    EventType
	Payload []byte
}

type EventType byte

const (
	EventOne EventType = 1
	EventMax EventType = 128
	New                = EventMax
	GetOne             = EventMax >> 1
	GetAll             = (EventMax >> 1) + EventMax

	TCostumer    = EventOne
	TAccount     = EventOne << 1
	TTransaction = (EventOne << 1) + EventOne

	NewAccountEvent     = New | TAccount
	NewCostumerEvent    = New | TCostumer
	NewTransactionEvent = New | TTransaction
)
