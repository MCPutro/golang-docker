package util

import (
	"github.com/MCPutro/golang-docker/internal/web/response"
	"github.com/gofiber/fiber/v2"
)

func WriteToResponseBody(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).
		JSON(response.Response{
			Status:  statusCode,
			Message: message,
			Data:    data,
		})
}
