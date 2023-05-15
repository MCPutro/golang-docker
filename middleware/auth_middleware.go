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
			return util.WriteToResponseBody(c, fiber.StatusUnauthorized, "unauthorized", nil)
		}

		//validation jwt
		validateToken, err := util.ValidateToken(strings.ReplaceAll(auth, "Bearer ", ""))
		if err != nil {
			return util.WriteToResponseBody(c, fiber.StatusUnauthorized, "invalid token. "+err.Error(), nil)
		}

		if validateToken.Valid {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusInternalServerError)

	}
}
