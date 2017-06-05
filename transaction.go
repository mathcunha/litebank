package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (t *Transaction) loadAccountFrom() error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	change := bson.M{"$rename": bson.M{"pending": "queue"}}
	if err := s.DB(database).C(t.From.collection()).Update(bson.M{"_id": t.From.Id}, change); err != nil {
		return err
	}
	if err := s.DB(database).C(t.From.collection()).Find(bson.M{"_id": t.From.Id}).One((&t.From)); err != nil {
		return err
	}

	//process queue
	for _, v := range t.From.Queue {
		t.From.Balance += v
	}

	return nil

}

func sameAccountTransaction(t *Transaction) (bool, error) {
	if err := t.loadAccountFrom(); err != nil {
		return false, err
	}

	invalid := t.Value < 0 && t.From.Balance < -t.Value

	if invalid {
		if len(t.From.Queue) > 0 {
			if err := updateAccountFrom((&t.From)); err != nil {
				return false, err
			}
		}
		return false, nil
	} else {
		t.From.Balance += t.Value
		if err := updateAccountFrom((&t.From)); err != nil {
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
	if err := t.loadAccountFrom(); err != nil {
		return false, err
	}

	if err := findOne(t.To.collection(), (&t.To), t.To.Id); err != nil {
		return false, err
	}

	invalid := t.From.Balance < t.Value

	if !invalid {
		t.From.Balance -= t.Value
		if err := updateAccountFrom((&t.From)); err != nil {
			return true, err
		}

		t.To.Queue = append(t.To.Queue, t.Value)
		if err := updateAccount(false, (&t.To)); err != nil {
			//try to restore the balance
			t.From.Balance += t.Value
			if err := updateAccountFrom((&t.From)); err != nil {
				return true, err
			}
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func updateAccountFrom(a *Account) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	change := bson.M{"$set": bson.M{"balance": a.Balance, "queue": []float64{}}}

	return s.DB(database).C(a.collection()).Update(bson.M{"_id": a.Id}, change)
}

func updateAccount(updateBalance bool, a *Account) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	change := bson.M{"$set": bson.M{"balance": a.Balance, "pending": a.Queue[0]}}
	if !updateBalance {
		change = bson.M{"$push": bson.M{"pending": a.Queue[0]}}
	}

	return s.DB(database).C(a.collection()).Update(bson.M{"_id": a.Id}, change)
}
