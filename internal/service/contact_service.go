package service

import (
	"fmt"
	"log"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"
	"github.com/MogboPython/belvaphilips_backend/pkg/utils"
)

// SendContactEmail sends a contact email to the admin
func SendContactEmail(req *model.ContactUsRequest) error {
	// Set up email
	to := config.Config("ADMIN_EMAIL")
	subject := fmt.Sprintf("Contact from: %s %s", req.Firstname, req.Lastname)
	body, err := utils.ParseTemplate("contact.html", req)
	if err != nil {
		log.Printf("failed to parse template: %v", err)
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Send email
	_, err = utils.SendEmail(to, subject, body)
	if err != nil {
		log.Printf("failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
