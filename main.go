package main

import (
	"ip-waf-helper/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
