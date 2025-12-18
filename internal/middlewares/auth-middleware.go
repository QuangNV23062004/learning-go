package middlewares

import (
	"learning-go/internal/config"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// helpers
func VerifyToken(tokenString string) (*jwt.Token, error) {
	issuer := config.GetEnv("JWT_ISSUER", "")
	secret := config.GetEnv("JWT_ACCESS_SECRET", "")
	token, err := jwt.ParseWithClaims(tokenString,
		jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

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

func SetClaimsToContext(c fiber.Ctx, token *jwt.Token) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("sub", claims["sub"])
		c.Locals("role", claims["role"])
	}
}

func AuthMiddleware(c fiber.Ctx) error {
	tokenString := ExtractToken(c)
	isPublic := c.Locals("public")

	//token not found: public => next, private => unauthorized
	if tokenString == "" {
		if isPublic == true {
			return c.Next()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid token",
		})
	}

	token, err := VerifyToken(tokenString)

	if err != nil || !token.Valid {
		if isPublic == true {
			return c.Next()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})

	}

	if token != nil {
		SetClaimsToContext(c, token)
	}

	return c.Next()
}
