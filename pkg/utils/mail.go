package utils

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strconv"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/gofiber/fiber/v2/log"
	gomail "gopkg.in/mail.v2"
)

func ParseTemplate(templateFileName string, data any) (string, error) {
	var body string

	templatePath := filepath.Join("templates", templateFileName)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}

	body = buf.String()

	return body, nil
}

func SendEmail(receiver, subject, body string) (bool, error) {
	// Create a new message
	message := gomail.NewMessage()
	message.SetHeader("From", config.Config("PLUNK_EMAIL"))
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Set up the SMTP dialer
	portStr := config.Config("MAIL_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Error("Invalid MAIL_PORT: %v", err)
		return false, err
	}

	dialer := gomail.NewDialer(config.Config("MAIL_HOST"), port, config.Config("PLUNK_USERNAME"), config.Config("PLUNK_API_KEY"))

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Error(err)
		return false, err
	}

	log.Info("Email sent successfully!")

	return true, nil
}
