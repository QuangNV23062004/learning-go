package http

import (
	"learning-go/internal/http"
	"learning-go/internal/pkg/products/application"
	"learning-go/internal/pkg/products/dtos"
	"learning-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type ProductHandler struct {
	Service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		Service: service,
	}
}

func (h *ProductHandler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	includeDeleted := c.Query("includeDeleted") == "true"
	role := c.Locals("role")
	if role == nil {
		role = ""
	}

	product, err := h.Service.GetProductByID(id, includeDeleted, role.(string))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		product, fiber.StatusOK,
	))
}

func (h *ProductHandler) GetProductsByUserID(c fiber.Ctx) error {
	userID := c.Params("id")
	includeDeleted := c.Query("includeDeleted") == "true"
	role := c.Locals("role")
	if role == nil {
		role = ""
	}
	products, err := h.Service.GetProductsByUserID(userID, includeDeleted, role.(string))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		products, fiber.StatusOK,
	))
}

func (h *ProductHandler) GetPaginatedProductsByUserID(c fiber.Ctx) error {
	var query dtos.PaginatedProductsQueryDto
	userID := c.Params("id")
	role := c.Locals("role")
	if role == nil {
		role = ""
	}
	if err := c.Bind().Query(&query); err != nil {
		err := http.ErrInvalidQuery
		return err
	}

	// Apply defaults
	query.ApplyDefaults()

	paginatedProducts, err := h.Service.GetPaginatedProductsByUserID(userID, query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted, role.(string))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		paginatedProducts, fiber.StatusOK,
	))
}

func (h *ProductHandler) GetAllProducts(c fiber.Ctx) error {
	includeDeleted := c.Query("includeDeleted") == "true"
	role := c.Locals("role")
	if role == nil {
		role = ""
	}
	products, err := h.Service.GetAllProducts(includeDeleted, role.(string))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		products, fiber.StatusOK,
	))

}

func (h *ProductHandler) GetPaginatedProducts(c fiber.Ctx) error {
	var query dtos.PaginatedProductsQueryDto
	role := c.Locals("role")
	if role == nil {
		role = ""
	}
	if err := c.Bind().Query(&query); err != nil {
		err := http.ErrInvalidQuery
		return err
	}

	// Apply defaults
	query.ApplyDefaults()

	paginatedProducts, err := h.Service.GetPaginatedProducts(query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted, role.(string))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		paginatedProducts, fiber.StatusOK,
	))

}

func (h *ProductHandler) CreateProduct(c fiber.Ctx) error {
	var body dtos.CreateProductDTO
	if err := c.Bind().Body(&body); err != nil {
		err := http.ErrInvalidBody
		return err
	}

	body.UserID = c.Locals("sub").(string)

	createdProduct, err := h.Service.CreateProduct(&body)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(utils.Success(
		createdProduct, fiber.StatusCreated,
	))
}

func (h *ProductHandler) UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var body dtos.UpdateProductDTO
	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	if err := c.Bind().Body(&body); err != nil {
		err := http.ErrInvalidBody
		return err
	}

	updatedProduct, err := h.Service.UpdateProduct(id, &body, sub, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		updatedProduct, fiber.StatusOK,
	))
}

func (h *ProductHandler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	deleted, err := h.Service.DeleteProduct(id, role, sub)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		deleted, fiber.StatusOK,
	))
}

func (h *ProductHandler) RestoreProduct(c fiber.Ctx) error {
	id := c.Params("id")
	restored, err := h.Service.RestoreProduct(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		restored, fiber.StatusOK,
	))
}
