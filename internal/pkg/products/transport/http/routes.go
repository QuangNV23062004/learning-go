package http

import (
	"learning-go/internal/middlewares"
	"learning-go/internal/pkg/users/enums"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handler    *ProductHandler
	jwtService *utils.JwtService
}

func NewRouter(handler *ProductHandler, jwtService *utils.JwtService) *Router {
	return &Router{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (r *Router) SetupRoutes(app fiber.Router) {
	product := app.Group("/products")

	product.Get("/",
		middlewares.MarkPublic(),
		middlewares.AuthMiddleware(r.jwtService),
		r.handler.GetPaginatedProducts)

	product.Get("/user/:id",
		middlewares.MarkPublic(),
		middlewares.AuthMiddleware(r.jwtService),
		r.handler.GetProductsByUserID)

	product.Get("/user/:id/paginated",
		middlewares.MarkPublic(),
		middlewares.AuthMiddleware(r.jwtService),
		r.handler.GetPaginatedProductsByUserID)

	product.Get("/all",
		middlewares.MarkPublic(),
		middlewares.AuthMiddleware(r.jwtService),
		r.handler.GetAllProducts)

	product.Get("/:id",
		middlewares.MarkPublic(),
		middlewares.AuthMiddleware(r.jwtService),
		r.handler.GetProductByID)

	product.Use(middlewares.AuthMiddleware(r.jwtService))

	product.Post("/",
		middlewares.RoleMiddleware(
			[]string{string(enums.Admin), string(enums.User)}),
		r.handler.CreateProduct)

	product.Patch("/:id",
		middlewares.RoleMiddleware(
			[]string{string(enums.Admin), string(enums.User)}),
		r.handler.UpdateProduct)

	product.Delete("/:id",
		middlewares.RoleMiddleware(
			[]string{string(enums.Admin), string(enums.User)}),
		r.handler.DeleteProduct)

	product.Post("/:id/restore",
		middlewares.RoleMiddleware(
			[]string{string(enums.Admin)}),
		r.handler.RestoreProduct)
}
