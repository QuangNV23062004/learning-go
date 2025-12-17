package http

import "github.com/gofiber/fiber/v3"

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
	auth.Post("/register", r.handler.Register)
	auth.Post("/login", r.handler.Login)
}
