package config

import (
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(userController controller.UserController) *fiber.App {
	router := fiber.New()

	//router.Get("/ping", func(ctx *fiber.Ctx) error {
	//	return ctx.SendString("Pong")
	//})
	//router.Post("/login", userController.Login)
	//
	//secureRoutes := router.Group("/", middleware.AuthMiddleware())
	//secureRoutes.Post("/registration", userController.Registration)
	//secureRoutes.Get("/user", middleware.AuthMiddleware(), userController.ShowAllUser)
	//secureRoutes.Get("/user/:uid", userController.ShowUser)
	//secureRoutes.Patch("/user/:uid", userController.UpdateUser)
	//secureRoutes.Delete("/user/:uid", userController.DeleteUser)

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
	secure.Patch("/user/:uid", userController.UpdateUser)
	secure.Delete("/user/:uid", userController.DeleteUser)

}
