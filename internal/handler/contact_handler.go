package handler

import (
	"github.com/MogboPython/belvaphilips_backend/internal/service"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// ContactUs handles the contact us form submissions
//
//	@Summary		Submit contact form
//	@Description	Submit contact form to notify admin
//	@Tags			contact
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.ContactUsRequest	true	"Contact information"
//	@Success		200		{object}	model.ResponseHTTP{}
//	@Failure		400		{object}	model.ResponseHTTP{}
//	@Failure		500		{object}	model.ResponseHTTP{}
//	@Router			/api/v1/contact [post]
func ContactUs(c *fiber.Ctx) error {
	var req model.ContactUsRequest

	v := validator.New()

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if err := v.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := service.SendContactEmail(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ResponseHTTP{
			Success: false,
			Message: "Failed to send contact email",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.ResponseHTTP{
		Success: true,
		Message: "Your message has been sent successfully",
		Data:    "",
	})
}
