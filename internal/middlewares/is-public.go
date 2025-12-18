package middlewares

import "github.com/gofiber/fiber/v3"

func MarkPublic(c fiber.Ctx) error {
	c.Locals("public", true)
	return c.Next()
}
