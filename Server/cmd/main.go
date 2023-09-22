package main

import (
	"example.com/application/middleware"
	"example.com/application/service"
	"example.com/config"
	infrastructure "example.com/infrastructure/persistence"
	_interface "example.com/interface/handler"
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
)

func main() {
	db, _ := config.NewDBConnection()

	userRepository := infrastructure.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := _interface.NewUserHandler(userService)
	middleware := middleware.NewMiddleware(userService)

	r := bunrouter.New()
	r.Use(middleware.RecoverMiddleware())
	r.Use(middleware.CorsMiddleware())

	r.POST("/user/create", userHandler.UserCreateHandle())
	r.POST("/user/get", userHandler.UserGetHandle())
	r.POST("/move", userHandler.MoveHandle())
	r.POST("/destroy", userHandler.DestroyHandle())
	r.GET("/users/get", userHandler.UserRankingGetHandle())

	log.Println("listening on http://localhost:8080")
	log.Println(http.ListenAndServe(":8080", r))
}
