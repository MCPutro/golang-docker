package util

import (
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/gofiber/fiber/v2"
)

func WriteToResponseBody(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).
		JSON(web.Response{
			Status:  statusCode,
			Message: message,
			Data:    data,
		})
}
