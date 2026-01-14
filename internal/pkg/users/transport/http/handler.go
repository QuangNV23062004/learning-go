package http

import (
	"encoding/json"

	httpError "github.com/QuangNV23062004/learning-go/internal/http"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/application"
	"github.com/QuangNV23062004/learning-go/internal/pkg/users/dtos"
	"github.com/QuangNV23062004/learning-go/internal/utils"

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
		err := httpError.ErrInvalidBody
		return err
	}

	user, err := h.service.Register(body)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success(
		"Email sent to "+user.Email, fiber.StatusCreated,
	))
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var body dtos.LoginDto
	if err := c.Bind().Body(&body); err != nil {
		err := httpError.ErrInvalidBody
		return err
	}

	credentials, err := h.service.Login(body)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(credentials.User)
	if err != nil {
		return err
	}

	c.Cookies("accessToken", credentials.AccessToken)
	c.Cookies("refreshToken", credentials.RefreshToken)
	c.Cookies("user", string(userJson))

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		credentials, fiber.StatusOK,
	))
}

func (h *UserHandler) VerifyUser(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		err := httpError.ErrUnauthorized
		return err
	}

	user, err := h.service.VerifyEmail(token)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		user, fiber.StatusOK,
	))
}

func (h *UserHandler) RefreshToken(c fiber.Ctx) error {

	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		err := httpError.ErrMissingRefreshToken
		return err
	}

	credentials, err := h.service.RefreshTokens(refreshToken)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(credentials.User)
	if err != nil {
		return err
	}

	c.Cookies("accessToken", credentials.AccessToken)
	c.Cookies("refreshToken", credentials.RefreshToken)
	c.Cookies("user", string(userJson))

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		credentials, fiber.StatusOK,
	))
}

// Users
func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
	includeDeleted := c.Query("includeDeleted", "false") == "true"
	users, err := h.service.GetAllUsers(includeDeleted)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		users, fiber.StatusOK,
	))
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")

	role := c.Locals("role").(string)
	includeDeleted := c.Query("includeDeleted", "false") == "true"
	user, err := h.service.GetUserByID(id, role, includeDeleted)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		user, fiber.StatusOK,
	))
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	sub := c.Locals("sub").(string)
	deleted, err := h.service.DeleteUser(id, role, sub)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		deleted, fiber.StatusOK,
	))
}

func (h *UserHandler) RestoreUser(c fiber.Ctx) error {
	id := c.Params("id")
	restored, err := h.service.RestoreUser(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		restored, fiber.StatusOK,
	))
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	sub := c.Locals("sub").(string)
	var body dtos.UpdateUserDto

	if err := c.Bind().Body(&body); err != nil {
		err := httpError.ErrInvalidBody
		return err
	}

	updatedUser, err := h.service.UpdateUser(id, body, sub)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success(
		updatedUser, fiber.StatusOK,
	))
}

func (h *UserHandler) PaginatedUsers(c fiber.Ctx) error {
	var query dtos.PaginatedUsersQueryDto
	if err := c.Bind().Query(&query); err != nil {
		err := httpError.ErrInvalidQuery
		return err
	}

	// Apply defaults
	query.ApplyDefaults()

	data, err := h.service.PaginatedUsers(query.Page, query.Limit, query.Search, query.SearchField, query.Order, query.SortBy, query.IncludeDeleted)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		data, fiber.StatusOK,
	))
}
