package main

import (
	_interface "example.com/interface/handler"
	"github.com/uptrace/bunrouter"
)

func main() {
	router := bunrouter.New()

	router.GET("/", _interface.BunHandle())
}
