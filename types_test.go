package main

import (
	"testing"
)

func TestEventTypes(t *testing.T) {
	t.Logf("New:%b, GetOne:%b, GetAll:%b, Costumer:%b, Account:%b, Transaction:%b", New, GetOne, GetAll, TCostumer, TAccount, TTransaction)
	t.Logf("New:%b, GetOne:%b, GetAll:%b, Costumer:%b, Account:%b, Transaction:%b", New, GetOne, GetAll, TCostumer<<3, TAccount<<3, TTransaction<<3)
	t.Logf("NewCostumer:%b, NewAccountEvent:%b, NewTransactionEvent:%b", NewCostumerEvent, NewAccountEvent, NewTransactionEvent)
}
