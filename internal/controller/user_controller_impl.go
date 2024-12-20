package controller

import (
	"errors"
	"fmt"
	user2 "github.com/MCPutro/golang-docker/internal/entity"
	"github.com/MCPutro/golang-docker/internal/service/user"
	"github.com/MCPutro/golang-docker/internal/util"
	"github.com/MCPutro/golang-docker/internal/util/logger"
	"github.com/MCPutro/golang-docker/internal/web/request"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type userControllerImpl struct {
	service user.Service
}

func NewUserController(service user.Service) UserController {
	return &userControllerImpl{service: service}
}

func (u *userControllerImpl) Login(c *fiber.Ctx) error {
	body := new(request.UserCreate)

	if err := c.BodyParser(&body); err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	// logging
	logger.ContextLogger(c.UserContext()).Debugf("%+v", body)

	userResponse, err := u.service.Login(c.UserContext(), body)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.WriteToResponseBody(c, fiber.StatusNotFound, "failed to login. "+err.Error(), nil)
		}

		return util.WriteToResponseBody(c, fiber.StatusUnauthorized, "failed to login. "+err.Error(), nil)
	}

	// logging
	logger.ContextLogger(c.UserContext()).Debugf("%+v", userResponse)

	//success message
	return util.WriteToResponseBody(c, fiber.StatusOK, "success", userResponse)
}

func (u *userControllerImpl) Registration(c *fiber.Ctx) error {
	body := new(request.UserCreate)

	if err := c.BodyParser(&body); err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	create, err := u.service.Registration(c.UserContext(), body)
	if err != nil {
		if errors.Is(err, util.ErrAlreadyUsed) {
			return util.WriteToResponseBody(c, fiber.StatusUnprocessableEntity, "failed to registration. "+err.Error(), nil)
		}

		return util.WriteToResponseBody(c, fiber.StatusInternalServerError, "failed to registration. "+err.Error(), nil)
	}

	//success message
	return util.WriteToResponseBody(c, fiber.StatusCreated, "success", create)
}

func (u *userControllerImpl) ShowAllUser(c *fiber.Ctx) error {
	get := c.Locals("UserId")
	fmt.Println("User id :", get)
	users, err := u.service.GetAll(c.UserContext())
	if err != nil {

		if errors.Is(err, util.ErrNotFound) {
			return util.WriteToResponseBody(c, fiber.StatusOK, "success", users)
		}

		return util.WriteToResponseBody(c, fiber.StatusInternalServerError, "failed to get users data. "+err.Error(), nil)
	}

	//success message
	return util.WriteToResponseBody(c, fiber.StatusOK, "success", users)
}

func (u *userControllerImpl) ShowUser(c *fiber.Ctx) error {
	suid := c.Params("uid", "-1")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, fmt.Sprintf("user id %s is not valid.", suid), nil)
	}

	getById, err := u.service.GetById(c.UserContext(), uid)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.WriteToResponseBody(c, fiber.StatusNotFound, err.Error(), nil)
		}

		return util.WriteToResponseBody(c, fiber.StatusInternalServerError, "failed to get user data. "+err.Error(), nil)

	}

	//success message
	return util.WriteToResponseBody(c, fiber.StatusOK, "success", getById)
}

func (u *userControllerImpl) UpdateUser(c *fiber.Ctx) error {
	body := new(user2.User)

	if err := c.BodyParser(&body); err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	suid := c.Params("uid", "")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, fmt.Sprintf("user id %s is not valid.", suid), nil)
	}

	body.Id = uid

	update, err := u.service.Update(c.UserContext(), body)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.WriteToResponseBody(c, fiber.StatusNotFound, "failed to update user. "+err.Error(), nil)
		}
		return util.WriteToResponseBody(c, fiber.StatusInternalServerError, "failed to update user. "+err.Error(), nil)
	}

	//success message
	return util.WriteToResponseBody(c, fiber.StatusOK, "success", update)
}

func (u *userControllerImpl) DeleteUser(c *fiber.Ctx) error {
	suid := c.Params("uid", "-1")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return util.WriteToResponseBody(c, fiber.StatusBadRequest, fmt.Sprintf("user id %s is not valid.", suid), nil)
	}

	err = u.service.Remove(c.UserContext(), uid)

	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return util.WriteToResponseBody(c, fiber.StatusNotFound, "failed to delete user. "+err.Error(), nil)
		}

		return util.WriteToResponseBody(c, fiber.StatusInternalServerError, "failed to delete user. "+err.Error(), nil)
	}

	//success message
	return util.WriteToResponseBody(c, fiber.StatusOK, "success", nil)
}
