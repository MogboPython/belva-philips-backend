package service

import (
	"fmt"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"

	"github.com/gofiber/fiber/v2/log"
)

// SendContactEmail sends a contact email to the admin
func SendContactEmail(req *model.ContactUsRequest) error {
	to := config.Config("ADMIN_EMAIL")
	subject := fmt.Sprintf("Contact from: %s %s", req.Firstname, req.Lastname)
	body, err := utils.ParseTemplate("contact.html", req)

	if err != nil {
		log.Error("failed to parse template: %v", err)
		return fmt.Errorf("failed to parse template: %w", err)
	}

	_, err = utils.SendEmail(to, subject, body)
	if err != nil {
		log.Error("failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
