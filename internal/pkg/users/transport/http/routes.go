package http

import (
	"learning-go/internal/middlewares"
	"learning-go/internal/pkg/users/enums"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handler *UserHandler
}

func NewRouter(handler *UserHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) SetupRoutes(app fiber.Router) {
	auth := app.Group("/auth")

	auth.Post("/register", middlewares.MarkPublic, r.handler.Register)

	auth.Post("/login", middlewares.MarkPublic, r.handler.Login)

	auth.Get("/verify", middlewares.MarkPublic, r.handler.VerifyUser)

	auth.Use(middlewares.AuthMiddleware)
	auth.Post("/refresh", r.handler.RefreshToken)

	user := app.Group("/users")

	user.Use(middlewares.AuthMiddleware)

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
