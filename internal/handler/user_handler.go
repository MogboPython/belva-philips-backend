package handler

import (
	"errors"
	"net/http"
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

// ResponseHTTP represents response body of this API
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserAccessToken creates a authorization token for an authenticated user
//
//	@Summary		Create a authorization token
//	@Description	Create a new authorization token with the provided information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.TokenRequestPayload	true	"User information"
//	@Success		201		{object}	ResponseHTTP{data=map[string]string}
//	@Failure		400		{object}	ResponseHTTP{}
//	@Failure		500		{object}	ResponseHTTP{}
//	@Router			/api/v1/user/token [post]
func (h *UserHandler) CreateUserAccessToken(c *fiber.Ctx) error {
	var payload model.TokenRequestPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	token, err := utils.GenerateToken(payload.UserSessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseHTTP{
			Success: false,
			Message: "Error generating OTP",
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
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
// @Security BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.CreateUserRequest	true	"User information"
//	@Success		201		{object}	ResponseHTTP{data=model.UserResponse}
//	@Failure		400		{object}	ResponseHTTP{}
//	@Failure		500		{object}	ResponseHTTP{}
//	@Router			/api/v1/user [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var payload model.CreateUserRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	user, err := h.userService.CreateUser(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
				Success: false,
				Message: "User with this email exists",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(http.StatusCreated).JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully registered user",
		Data:    *user,
	})
}

// GetAllUsers is a function to get all user data from the database
//
//	@Summary		Get all users
//	@Description	Fetch a paginated list of users from the database
//	@Tags			users
//
// @Security BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (default is 1)"
//	@Param			limit	query		int	false	"Number of users per page (default is 10)"
//	@Success		200		{array}		ResponseHTTP{data=[]model.UserResponse}
//	@Failure		500		{object}	ResponseHTTP{}
//	@Router			/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")    // Default to page 1
	limitStr := c.Query("limit", "10") // Default to limit 10

	users, err := h.userService.GetAllUsers(pageStr, limitStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(http.StatusCreated).JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved users.",
		Data:    users,
	})
}

// GetUserByID is a function to get a user by ID
//
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			users
//
// @Security BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	ResponseHTTP{data=model.UserResponse}
//	@Failure		404	{object}	ResponseHTTP{}
//	@Failure		500	{object}	ResponseHTTP{}
//	@Router			/api/v1/user/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ResponseHTTP{
				Success: false,
				Message: "User not found",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		}) // return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusCreated).JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully found user.",
		Data:    *user,
	})
}

// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
