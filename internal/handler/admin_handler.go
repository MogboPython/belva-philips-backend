package handler

import (
	"errors"
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// AdminLogin is a handler for creating an authorization token for an admin user
//
//	@Summary		Logs admin user into the system
//	@Description	Create a new authorization token with the provided information
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			username	query		string	true	"The user name for login"
//	@Param			password	query		string	true	"The password name for login"
//	@Success		201			{object}	model.ResponseHTTP{data=map[string]string}
//	@Failure		400			{object}	model.ResponseHTTP{}
//	@Failure		500			{object}	model.ResponseHTTP{}
//	@Router			/api/v1/admin/login [post]
func (h *AdminHandler) AdminLogin(c *fiber.Ctx) error {
	var payload model.AdminLoginRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
		})
	}

	token, err := h.adminService.Login(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "incorrect username or password") {
			return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
				Success: false,
				Message: "Incorrect username or password",
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
		Message: "Success get access token",
		Data: map[string]string{
			"access_token": token,
		},
	})
}

// GetAllUsers is a function to get all user data from the database
//
//	@Summary		Get all users
//	@Description	Fetch a paginated list of users from the database
//	@Tags			admin
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (default is 1)"
//	@Param			limit	query		int	false	"Number of users per page (default is 10)"
//	@Success		200		{array}		model.ResponseHTTP{data=[]model.UserResponse}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/admin/get_users [get]
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")    // Default to page 1
	limitStr := c.Query("limit", "10") // Default to limit 10

	users, err := h.adminService.GetAllUsers(pageStr, limitStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Successfully retrieved users.",
		Data:    users,
	})
}

// GetUserByID is a function to get a user by ID
//
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			admin
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	model.ResponseHTTP{data=model.UserResponse}
//	@Failure		404	{object}	model.ResponseHTTP{}
//	@Failure		500	{object}	model.ResponseHTTP{}
//	@Router			/api/v1/user/{id} [get]
func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.adminService.GetUserByID(id)
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
