package http

import (
	"github.com/QuangNV23062004/learning-go/internal/middlewares"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/enums"
	"github.com/QuangNV23062004/learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handler    *UserHandler
	jwtService *utils.JwtService
}

func NewRouter(handler *UserHandler, jwtService *utils.JwtService) *Router {
	return &Router{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (r *Router) SetupRoutes(app fiber.Router) {
	auth := app.Group("/auth")

	auth.Post("/register", middlewares.MarkPublic(), r.handler.Register)

	auth.Post("/login", middlewares.MarkPublic(), r.handler.Login)

	auth.Get("/verify", middlewares.MarkPublic(), r.handler.VerifyUser)

	auth.Post("/refresh", middlewares.MarkPublic(), r.handler.RefreshToken)

	user := app.Group("/users")

	user.Use(middlewares.AuthMiddleware(r.jwtService))

	user.Get("/",
		middlewares.RoleMiddleware([]string{string(enums.Admin)}),
		r.handler.PaginatedUsers)

	user.Get("/all",
		middlewares.RoleMiddleware([]string{string(enums.Admin)}),
		r.handler.GetAllUsers)

	user.Get("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin),
			string(enums.User)}), r.handler.GetUserByID)

	user.Patch("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.UpdateUser)

	user.Delete("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.DeleteUser)

	user.Post("/:id/restore",
		middlewares.RoleMiddleware([]string{string(enums.Admin)}),
		r.handler.RestoreUser)

}
