package handler

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// TODO: remove this from swagger
// CreateUserAccessToken creates a authorization token for an authenticated user
//
//	@Summary		Create a authorization token
//	@Description	Create a new authorization token with the provided information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.TokenRequestPayload	true	"User information"
//	@Success		201		{object}	model.ResponseHTTP{data=map[string]string}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/user/token [post]
func (h *UserHandler) CreateUserAccessToken(c *fiber.Ctx) error {
	var payload model.TokenRequestPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	token, err := utils.GenerateToken(payload.UserSessionID)
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
//	@Param			user	body		model.CreateUserRequest	true	"User information"
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

// func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	user, err := h.userService.GetUserByID(id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return c.Status(fiber.StatusNotFound).JSON(model.ResponseHTTP{
// 				Success: false,
// 				Message: "User not found",
// 				Data:    nil,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
// 			Success: false,
// 			Message: "Internal server error",
// 			Data:    nil,
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
// 		Success: true,
// 		Message: "Successfully found user.",
// 		Data:    *user,
// 	})
// }

// GetUserByEmail is a function to get a user by Email
//
//	@Summary		Get user by Email
//	@Description	Get user by Email
//	@Tags			users
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			email	query		string	true	"User Email"
//	@Success		200		{object}	model.ResponseHTTP{data=model.UserResponse}
//	@Failure		404		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/users/get_user [get]
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	var payload model.GetUserByEmailRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	user, err := h.userService.GetUserByEmail(&payload)
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
