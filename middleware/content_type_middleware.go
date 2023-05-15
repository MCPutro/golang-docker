package middleware

import (
	"fmt"
	"github.com/MCPutro/golang-docker/util"
	"github.com/gofiber/fiber/v2"
)

func ContentTypeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut {
			contentType := c.Get(fiber.HeaderContentType, "")
			if contentType == "" {
				c.Request().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
				return c.Next()
			}

			if contentType == fiber.MIMEApplicationJSON || contentType == fiber.MIMEApplicationJSONCharsetUTF8 {
				return c.Next()
			} else {
				return util.WriteToResponseBody(c, fiber.StatusUnsupportedMediaType, fmt.Sprintf("Unsupported Content Type [%s].", contentType), nil)
			}

		} else {
			return c.Next()
		}
	}
}
