package main

import (
	"github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/database"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/MCPutro/golang-docker/util/logger"
	"log"
)

func main() {
	loggers := logger.NewLogger()
	loggers.Infoln("Application Starting")

	db, err := database.InitDatabase()
	if err != nil {
		loggers.Errorf("Database initialization failed with error: %s", err)
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

	loggers.Infof("Application listening on port %s", PORT)

	err = router.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(err)
	}
}
