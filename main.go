package main

import (
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/database"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New()

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Pong")
	})
	app.Post("/login", userController.Login)
	app.Post("/registration", userController.Registration)
	app.Get("/user", userController.ShowAllUser)
	app.Get("/user/:uid", userController.ShowUser)
	app.Patch("/user/:uid", userController.UpdateUser)
	app.Delete("/user/:uid", userController.DeleteUser)

	PORT := "9999"
	if PORT == "" {
		PORT = "9999"
	}

	log.Println("Running in port", PORT)

	err = app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(err)
	}
}
