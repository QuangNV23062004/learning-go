package http

import (
	"learning-go/internal/middlewares"
	"learning-go/internal/pkg/users/enums"
	"learning-go/internal/utils"

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
		middlewares.RoleMiddleware([]string{string(enums.Admin), string(enums.User)}))

	ordersGroup.Get("/:id",
		r.handler.FindOrderByID)

	ordersGroup.Get("/user/:id",
		r.handler.FindOrdersByUserID)

	ordersGroup.Get("/user/:id/paginated",
		r.handler.FindPaginatedOrdersByUserID)

	ordersGroup.Post("/",
		r.handler.CreateOrder)

	ordersGroup.Put("/:id",
		r.handler.UpdateOrder)

	ordersGroup.Delete("/:id",
		r.handler.DeleteOrder)

	ordersGroup.Use(
		middlewares.RoleMiddleware([]string{string(enums.Admin)}))

	ordersGroup.Get("/", r.handler.PaginatedOrders)
	ordersGroup.Get("/all", r.handler.FindAllOrders)

}
