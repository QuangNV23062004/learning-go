package middlewares

import "github.com/gofiber/fiber/v3"

func MarkPublic() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("public", true)
		return c.Next()
	}
}
