package main

import (
	"example.com/application/service"
	"example.com/config"
	infrastructure "example.com/infrastructure/persistence"
	_interface "example.com/interface/handler"
	"github.com/uptrace/bunrouter"
)

func main() {
	db, _ := config.NewDBConnection()

	userRepository := infrastructure.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := _interface.NewUserHandler(userService)

	router := bunrouter.New()
	router.POST("/login", userHandler.LoginHandle())
	router.POST("/move", userHandler.MoveHandle())
	router.POST("/destroy", userHandler.DestroyHandle())
	router.GET("/other", userHandler.OtherPlayerHandle())
}
