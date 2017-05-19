package main

import (
	"flag"
)

var server bool
var consumer bool

func init() {
	flag.BoolVar(&server, "server", false, "litebank starts as http server")
	flag.BoolVar(&consumer, "consumer", true, "litebank starts as event consumer")
}

func main() {
	flag.Parse()
	if server {
		err := listen()
		if err != nil {
			panic(err)
		}
	}
	if consumer {
		err := consume()
		if err != nil {
			panic(err)
		}
	}
}
