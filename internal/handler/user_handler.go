package handler

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/MogboPython/belvaphilips_backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) CreateUserAccessToken(c *fiber.Ctx) error {
	var payload model.TokenRequestPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	token, err := utils.GenerateToken(payload.UserSessionID, "authenticated")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Error generating access token",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Success get access token",
		Data: map[string]string{
			"access_token": token,
		},
	})
}

// CreateUser creates a new user
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided information
//	@Tags			users
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateUserRequest	true	"User information"
//	@Success		201		{object}	model.ResponseHTTP{data=model.UserResponse}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var payload model.CreateUserRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	user, err := h.userService.CreateUser(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: "User with this email exists",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully registered user",
		Data:    *user,
	})
}

// GetUserByID is a function to get a user by ID
//
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			users
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	model.ResponseHTTP{data=model.UserResponse}
//	@Failure		404	{object}	model.ResponseHTTP{}
//	@Failure		500	{object}	model.ResponseHTTP{}
//	@Router			/api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "User not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully found user.",
		Data:    *user,
	})
}

// UpdateMembershipStatus is a function to update a user's membership status
//
//	@Summary		Update the membership status of a user
//	@Description	Update the membership status of a user
//	@Tags			users
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"User ID"
//	@Param			request	body		model.OrderStatusChangeRequest	true	"Status update"
//	@Success		200		{object}	model.ResponseHTTP{data=model.MembershipStatusChangeRequest}
//	@Failure		404		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/users/{id}/membership [put]
func (h *UserHandler) UpdateMembershipStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload model.MembershipStatusChangeRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	user, err := h.userService.UpdateUserMembershipStatusChange(id, &payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
				Success: false,
				Message: "User not found",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully updated user membership status",
		Data:    *user,
	})
}
