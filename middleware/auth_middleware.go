package middleware

import (
	"github.com/MCPutro/golang-docker/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

//func for authentication token jwt

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get(fiber.HeaderAuthorization, "")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		//validation jwt
		validateToken, err := util.ValidateToken(strings.ReplaceAll(auth, "Bearer ", ""))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())

		}

		if validateToken.Valid {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusInternalServerError)

	}
}
