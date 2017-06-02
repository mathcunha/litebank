package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	Id    string `bson:"_id"`
	From  Account
	To    Account
	Value float64
}

func sameAccountTransaction(t *Transaction) (bool, error) {
	//load the account
	if err := findOne(t.From.collection(), (&t.From), t.From.Id); err != nil {
		return false, err
	}

	//process queue
	for _, v := range t.From.Queue {
		t.From.Balance += v
	}
	t.From.Queue = t.From.Queue[0:0]

	invalid := t.Value < 0 && t.From.Balance < -t.Value

	if !invalid {
		t.From.Balance += t.Value
		if err := updateAccount(true, (&t.From)); err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}

func processTransaction(t *Transaction) (bool, error) {
	if t.From.Id == t.To.Id {
		return sameAccountTransaction(t)
	}

	//load the accounts
	if err := findOne(t.From.collection(), (&t.From), t.From.Id); err != nil {
		return false, err
	}
	if err := findOne(t.To.collection(), (&t.To), t.To.Id); err != nil {
		return false, err
	}

	//process queue
	for _, v := range t.From.Queue {
		t.From.Balance += v
	}
	t.From.Queue = t.From.Queue[0:0]

	invalid := t.From.Balance < t.Value

	if !invalid {
		t.From.Balance -= t.Value
		if err := updateAccount(true, (&t.From)); err != nil {
			return true, err
		}
		t.To.Queue = append(t.To.Queue, t.Value)
		if err := updateAccount(false, (&t.To)); err != nil {
			//try to restore the balance
			t.From.Balance += t.Value
			if err := updateAccount(true, (&t.From)); err != nil {
				return true, err
			}
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func updateAccount(updateBalance bool, a *Account) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	change := bson.M{"$set": bson.M{"balance": a.Balance, "queue": a.Queue}}
	if !updateBalance {
		change = bson.M{"$set": bson.M{"queue": a.Queue}}
	}

	return s.DB(database).C(a.collection()).Update(bson.M{"_id": a.Id}, change)
}