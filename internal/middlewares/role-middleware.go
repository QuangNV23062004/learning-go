package middlewares

import (
	"github.com/QuangNV23062004/learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

func RoleMiddleware(requiredRoles []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(
				utils.Error("Access denied: insufficient permissions", fiber.StatusForbidden))
		}

		hasRole := false
		for _, r := range requiredRoles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(
				utils.Error("Access denied: insufficient permissions", fiber.StatusForbidden))
		}

		return c.Next()
	}
}
