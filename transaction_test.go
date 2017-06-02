package main

import (
	"testing"
)

func TestProcessTransaction(test *testing.T) {
	//t := Transaction{Value: 300, From: Account{Id: "2f8c3dd7-61cd-4113-868d-db21bb40890b"}, To: Account{Id: "9eb8207f-09ac-4fcf-b2c4-a1bed55706c1"}}
	t := Transaction{Value: 200, To: Account{Id: "2f8c3dd7-61cd-4113-868d-db21bb40890b"}, From: Account{Id: "9eb8207f-09ac-4fcf-b2c4-a1bed55706c1"}}
	//t := Transaction{Value: 200, To: Account{Id: "9eb8207f-09ac-4fcf-b2c4-a1bed55706c1"}, From: Account{Id: "9eb8207f-09ac-4fcf-b2c4-a1bed55706c1"}}
	valid, err := processTransaction(&t)
	test.Logf("%v, %v", valid, err)
}
