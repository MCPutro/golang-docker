package main

import (
	"github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/database"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"log"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("failed create database connection, error : %s", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	PORT := config.App_Port
	if PORT == "" {
		PORT = "1234"
	}

	router := config.NewRouter(userController)

	log.Println("Running in port", PORT)

	err = router.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(err)
	}
}
