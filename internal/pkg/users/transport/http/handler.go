package http

import (
	"encoding/json"
	"learning-go/internal/pkg/users/application"
	"learning-go/internal/pkg/users/dtos"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Auth
func (h *UserHandler) Register(c fiber.Ctx) error {
	var body dtos.RegisterDto
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.service.Register(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"Data": "Email sent to " + user.Email})
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var body dtos.LoginDto
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	credentials, err := h.service.Login(body)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userJson, err := json.Marshal(credentials.User)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookies("accessToken", credentials.AccessToken)
	c.Cookies("refreshToken", credentials.RefreshToken)
	c.Cookies("user", string(userJson))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": credentials,
	})
}

func (h *UserHandler) VerifyUser(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing token",
		})
	}

	user, err := h.service.VerifyEmail(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data":    user,
		"message": "Email verified successfully for user " + user.Email,
	})
}

func (h *UserHandler) RefreshToken(c fiber.Ctx) error {

	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing refresh token",
		})
	}

	credentials, err := h.service.RefreshTokens(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userJson, err := json.Marshal(credentials.User)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookies("accessToken", credentials.AccessToken)
	c.Cookies("refreshToken", credentials.RefreshToken)
	c.Cookies("user", string(userJson))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": credentials,
	})
}

// Users
func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
	includeDeleted := c.Query("includeDeleted", "false") == "true"
	users, err := h.service.GetAllUsers(includeDeleted)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": users,
	})
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	includeDeleted := c.Query("includeDeleted", "false") == "true"
	user, err := h.service.GetUserByID(id, role, includeDeleted)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": user,
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	deleted, err := h.service.DeleteUser(id, role, sub)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": deleted,
	})
}

func (h *UserHandler) RestoreUser(c fiber.Ctx) error {
	id := c.Params("id")
	restored, err := h.service.RestoreUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": restored,
	})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	sub := c.Locals("sub").(string)
	var body dtos.UpdateUserDto

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedUser, err := h.service.UpdateUser(id, body, sub)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Data": updatedUser,
	})
}

func (h *UserHandler) PaginatedUsers(c fiber.Ctx) error {
	var query dtos.PaginatedUsersQueryDto
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Apply defaults
	query.ApplyDefaults()

	data, err := h.service.PaginatedUsers(query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
