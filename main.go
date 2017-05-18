package main

import (
	"flag"
)

var server bool

func init() {
	flag.BoolVar(&server, "server", false, "starts as litebank http server")
}

func main() {
	flag.Parse()
	if server {
		err := listen()
		if err != nil {
			panic(err)
		}
	}
}
