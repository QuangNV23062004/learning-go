package middlewares

import "github.com/gofiber/fiber/v3"

func RoleMiddleware(requiredRoles []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: role not found",
			})
		}

		hasRole := false
		for _, r := range requiredRoles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions",
			})
		}

		return c.Next()
	}
}
