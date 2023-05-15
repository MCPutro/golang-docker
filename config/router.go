package config

import (
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(userController controller.UserController) *fiber.App {
	router := fiber.New()

	router.Use(middleware.ContentTypeMiddleware())

	guestRoutes(router, userController)
	secureRoutes(router, userController)

	return router
}

func guestRoutes(router *fiber.App, userController controller.UserController) {
	router.Get("/ping", func(ctx *fiber.Ctx) error { return ctx.SendString("Pong") })
	router.Post("/login", userController.Login)

}
func secureRoutes(router *fiber.App, userController controller.UserController) {
	secure := router.Group("/", middleware.AuthMiddleware())

	secure.Post("/registration", userController.Registration)
	secure.Get("/user", userController.ShowAllUser)
	secure.Get("/user/:uid", userController.ShowUser)
	secure.Put("/user/:uid", userController.UpdateUser)
	secure.Delete("/user/:uid", userController.DeleteUser)

}
