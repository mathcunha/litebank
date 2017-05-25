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
	Id           string
	Name         string
	Creation     time.Time
	Accounts     []Account
	Transactions []Transaction
}

func (c *Costumer) collection() string {
	return "costumer"
}

type Account struct {
	Id       string
	Number   string
	Balance  float64
	Costumer Costumer
}

func (a *Account) collection() string {
	return "account"
}

type Transaction struct {
	Id    string
	From  Costumer
	To    Costumer
	Value float64
}

func (c *Transaction) collection() string {
	return "transaction"
}

type Event struct {
	Type EventType
	Body string
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
