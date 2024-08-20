package main

import (
	"github.com/48Club/ip-waf-helper/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
