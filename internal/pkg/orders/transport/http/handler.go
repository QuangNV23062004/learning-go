package http

import (
	"learning-go/internal/http"
	"learning-go/internal/pkg/orders/application"
	"learning-go/internal/pkg/orders/dtos"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	service *application.OrderService
}

func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) FindOrderByID(c fiber.Ctx) error {

	role := c.Locals("role").(string)
	id := c.Params("id")
	includeDeleted := c.Query("includeDeleted") == "true"

	order, err := h.service.FindOrderByID(id, includeDeleted, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(order, fiber.StatusOK))
}

func (h *OrderHandler) FindOrdersByUserID(c fiber.Ctx) error {

	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	userID := c.Params("id")
	includeDeleted := c.Query("includeDeleted") == "true"
	orders, err := h.service.FindOrdersByUserID(userID, includeDeleted, sub, role)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(orders, fiber.StatusOK))
}

func (h *OrderHandler) FindPaginatedOrdersByUserID(c fiber.Ctx) error {

	var query dtos.PaginatedProductsQueryDto

	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	userID := c.Params("id")

	if err := c.Bind().Query(&query); err != nil {
		err := http.ErrInvalidQuery
		return err
	}

	// Apply defaults
	query.ApplyDefaults()

	orders, err := h.service.FindPaginatedOrdersByUserID(userID, query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted, sub, role)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(orders, fiber.StatusOK))
}

func (h *OrderHandler) FindAllOrders(c fiber.Ctx) error {
	role := c.Locals("role").(string)
	includeDeleted := c.Query("includeDeleted") == "true"

	orders, err := h.service.FindAllOrders(includeDeleted, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(orders, fiber.StatusOK))
}

func (h *OrderHandler) PaginatedOrders(c fiber.Ctx) error {

	var query dtos.PaginatedProductsQueryDto
	role := c.Locals("role").(string)
	if err := c.Bind().Query(&query); err != nil {
		err := http.ErrInvalidQuery
		return err
	}

	// Apply defaults
	query.ApplyDefaults()

	orders, err := h.service.Paginated(query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(orders, fiber.StatusOK))
}

func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
	var orderDto dtos.CreateOrderDTO
	sub := c.Locals("sub").(string)
	if err := c.Bind().Body(&orderDto); err != nil {
		return http.ErrInvalidBody
	}

	order, err := h.service.Create(&orderDto, sub)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success(order, fiber.StatusCreated))
}

func (h *OrderHandler) UpdateOrder(c fiber.Ctx) error {
	var orderDto dtos.UpdateOrderDTO
	id := c.Params("id")
	sub := c.Locals("sub").(string)
	role := c.Locals("role").(string)

	if err := c.Bind().Body(&orderDto); err != nil {
		return http.ErrInvalidBody
	}

	updatedOrder, err := h.service.Update(id, &orderDto, sub, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(updatedOrder, fiber.StatusOK))

}

func (h *OrderHandler) DeleteOrder(c fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)

	deleted, err := h.service.Delete(id, sub, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(deleted, fiber.StatusOK))
}
