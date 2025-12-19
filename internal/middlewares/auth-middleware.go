package middlewares

import (
	"learning-go/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractToken(c fiber.Ctx) string {
	token := c.Cookies("accessToken")
	if token == "" {

		authorization := c.Get("Authorization")
		authorizationParts := strings.Split(authorization, " ")

		if len(authorizationParts) == 2 && strings.ToLower(authorizationParts[0]) == "bearer" {
			token = authorizationParts[1]
		} else {
			token = ""
		}
	}

	return token
}

func SetClaimsToContext(c fiber.Ctx, claims jwt.MapClaims) {
	c.Locals("sub", claims["sub"])
	c.Locals("role", claims["role"])
}

func AuthMiddleware(jwtService *utils.JwtService) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := ExtractToken(c)
		isPublic := c.Locals("public").(bool)

		//token not found: public => next, private => unauthorized
		if tokenString == "" {
			if isPublic == true {
				return c.Next()
			}
			return c.Status(fiber.StatusUnauthorized).JSON(
				utils.Error("Missing or invalid token", fiber.StatusUnauthorized),
			)
		}

		//without checking token valid
		claims, err := jwtService.VerifyAccessToken(tokenString)

		if err != nil {
			if isPublic == true {
				return c.Next()
			}
			return c.Status(fiber.StatusUnauthorized).JSON(
				utils.Error("Invalid token", fiber.StatusUnauthorized),
			)

		}

		if claims != nil {
			SetClaimsToContext(c, claims)
		}

		return c.Next()
	}
}
