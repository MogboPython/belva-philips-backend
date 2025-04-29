package handler

import (
	"strings"

	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService service.AdminService
	validator    *validator.Validator
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		validator:    validator.New(),
	}
}

// AdminLogin is a handler for creating an authorization token for an admin user
//
//	@Summary		Logs admin user into the system
//	@Description	Create a new authorization token with the provided information
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.AdminLoginRequest	true	"Login information"
//	@Success		201		{object}	model.ResponseHTTP{data=map[string]string}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
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

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
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
