package http

import (
	"github.com/QuangNV23062004/learning-go/internal/middlewares"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/enums"
	"github.com/QuangNV23062004/learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handler    *OrderHandler
	jwtService *utils.JwtService
}

func NewRouter(handler *OrderHandler, jwtService *utils.JwtService) *Router {
	return &Router{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (r *Router) SetupRoutes(appGroup fiber.Router) {
	ordersGroup := appGroup.Group("/orders")
	ordersGroup.Use(
		middlewares.AuthMiddleware(r.jwtService),
	)

	ordersGroup.Get("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.FindOrderByID)

	ordersGroup.Get("/user/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.FindOrdersByUserID)

	ordersGroup.Get("/user/:id/paginated",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.FindPaginatedOrdersByUserID)

	ordersGroup.Post("/",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.CreateOrder)

	ordersGroup.Patch("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.UpdateOrder)

	ordersGroup.Delete("/:id",
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}),
		r.handler.DeleteOrder)

	ordersGroup.Use(
		middlewares.RoleMiddleware([]string{string(enums.Admin)}))

	ordersGroup.Get("/", r.handler.PaginatedOrders)
	ordersGroup.Get("/all", r.handler.FindAllOrders)

}
