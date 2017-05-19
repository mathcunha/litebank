package main

import (
	"time"
)

type Entity interface {
	create() (*Event, error)
	newEntity() Entity
}

type Costumer struct {
	Id           string
	Name         string
	Creation     time.Time
	Accounts     []Account
	Transactions []Transaction
}

type Account struct {
	Id       string
	Number   string
	Balance  float64
	Costumer Costumer
}

type Transaction struct {
	Id    string
	From  Costumer
	To    Costumer
	Value float64
}

type Event struct {
	Type EventType
	Body string
}

type EventType byte

var NewCostumerEvent EventType = 0
var NewAccountEvent EventType = 10
var NewTransactionEvent EventType = 20
