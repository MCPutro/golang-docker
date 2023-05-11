package controller

import (
	"fmt"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type userControllerImpl struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userControllerImpl{service: service}
}

func (u *userControllerImpl) Login(c *fiber.Ctx) error {
	body := new(web.UserCreateRequest)

	if err := c.BodyParser(&body); err != nil {
		return c.SendString(err.Error())
	}

	user, err := u.service.Login(c.UserContext(), body)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(user)
}

func (u *userControllerImpl) Registration(c *fiber.Ctx) error {
	body := new(web.UserCreateRequest)

	if err := c.BodyParser(&body); err != nil {
		return c.SendString(err.Error())
	}

	create, err := u.service.Registration(c.UserContext(), body)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(create)
}

func (u *userControllerImpl) ShowAllUser(c *fiber.Ctx) error {
	users, err := u.service.GetAll(c.UserContext())
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(users)
}

func (u *userControllerImpl) ShowUser(c *fiber.Ctx) error {
	suid := c.Params("uid", "-1")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.SendString(fmt.Sprintf("User ID %s not valid", suid))
	}

	user, err := u.service.GetById(c.UserContext(), uid)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(user)
}

func (u *userControllerImpl) UpdateUser(c *fiber.Ctx) error {
	suid := c.Params("uid", "-1")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.SendString(fmt.Sprintf("User ID %s not valid", suid))
	}

	body := new(model.User)

	if err := c.BodyParser(&body); err != nil {
		return c.SendString(err.Error())
	}
	body.Id = uid

	update, err := u.service.Update(c.UserContext(), body)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(update)

}

func (u *userControllerImpl) DeleteUser(c *fiber.Ctx) error {
	suid := c.Params("uid", "-1")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.SendString(fmt.Sprintf("User ID %s not valid", suid))
	}

	err = u.service.Remove(c.UserContext(), uid)

	if err != nil {
		return c.SendString(err.Error())
	}

	return c.SendString("berhasil hapus user dengan id : " + suid)
}
