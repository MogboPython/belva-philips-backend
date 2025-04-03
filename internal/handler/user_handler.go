package handler

import (
	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func CreateUser(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// CreateUserAccessToken creates a authorization token for an authenticated user
// @Summary Create a authorization token
// @Description Create a new authorization token with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.TokenRequestPayload true "User information"
// @Success 201 {string} string
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
func (h *UserHandler) CreateUserAccessToken(c *fiber.Ctx) error {
	var payload model.TokenRequestPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	token, err := utils.GenerateToken(payload.UserSessionID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"status":       "success",
		"access_token": token,
	})
}
