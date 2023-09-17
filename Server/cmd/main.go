package main

import (
	"example.com/application/service"
	"example.com/config"
	db_init "example.com/db/init"
	infrastructure "example.com/infrastructure/persistence"
	_interface "example.com/interface/handler"
	"github.com/uptrace/bunrouter"
)

func main() {
	db, _ := config.NewDBConnection()

	db_init.CreateTable(db)

	userRepository := infrastructure.NewUserRepository(db)

	service.NewUserService(userRepository)

	router := bunrouter.New()
	router.POST("/login", _interface.LoginHandle())
	router.POST("/move", _interface.MoveHandle())
	router.POST("/destroy", _interface.DestroyHandle())
	router.GET("/other", _interface.OtherPlayerHandle())
}
