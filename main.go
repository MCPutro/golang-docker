package main

import (
	"github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/internal/controller"
	"github.com/MCPutro/golang-docker/internal/database"
	"github.com/MCPutro/golang-docker/internal/repository"
	"github.com/MCPutro/golang-docker/internal/routes"
	"github.com/MCPutro/golang-docker/internal/service"
	"github.com/MCPutro/golang-docker/internal/util/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	loggers := logger.NewLogger(logrus.DebugLevel)
	loggers.Infoln("Application Starting")

	config.NewConfig()

	db, err := database.InitDatabase()
	if err != nil {
		loggers.Errorf("Database initialization failed with error: %s", err)
	}
	defer db.Close()

	repoManager := repository.NewRepositoryManager()
	serviceManager := service.NewServiceManager(repoManager, db)
	userController := controller.NewUserController(serviceManager.UserService())

	PORT := config.App_Port
	if PORT == "" {
		PORT = "1234"
	}

	router := routes.NewRouter(userController)

	loggers.Infof("Application listening on port %s", PORT)

	err = router.Listen(":" + PORT)
	if err != nil {
		// log.Fatalln(err)
		loggers.Errorln("Application listening failed with error :", err)
	}
}
