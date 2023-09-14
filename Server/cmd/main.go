package main

import (
	_interface "example.com/interface/handler"
	"github.com/uptrace/bunrouter"
)

func main() {
	router := bunrouter.New()

	router.POST("/login", _interface.LoginHandle())
	router.POST("/move", _interface.MoveHandle())
	router.POST("/destroy", _interface.DestroyHandle())
	router.GET("/other", _interface.OtherPlayerHandle())
}
