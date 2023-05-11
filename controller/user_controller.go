package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	Login(c *fiber.Ctx) error
	Registration(c *fiber.Ctx) error
	ShowAllUser(c *fiber.Ctx) error
	ShowUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}
