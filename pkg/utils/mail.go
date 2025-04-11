package utils

// import (
// 	"bytes"
// 	"html/template"
// 	"net/smtp"

// 	"github.com/MogboPython/belvaphilips_backend/internal/config"
// )

// Request struct
// type Request struct {
// 	from    string
// 	to      []string
// 	subject string
// 	body    string
// }

// func NewRequest(to []string, subject, body string) *Request {
// 	return &Request{
// 		to:      to,
// 		subject: subject,
// 		body:    body,
// 	}
// }

// type TemplateData struct {
// 	Name    string
// 	Email   string
// 	Message string
// }

// var auth smtp.Auth

// func SendEmail(reciever, name, email, message string) (bool, error) {
// 	auth = smtp.PlainAuth(config.Config("PLUNK_USERNAME"), config.Config("PLUNK_EMAIL"), config.Config("PLUNK_API_KEY"), "smtp.useplunk.com")

// 	templateData := &TemplateData{
// 		Name:    name,
// 		Email:   email,
// 		Message: message,
// 	}

// 	email_req := NewRequest([]string{"junk@junk.com"}, "Hello Junk!", "Hello, World!")
// 	if err := parseTemplate("template.html", templateData); err != nil {
// 		return false, err
// 	}

// 	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
// 	subject := "Subject: New Contact Information from Website!\n"
// 	msg := []byte(subject + mime + "\n" + r.body)
// 	addr := "smtp.useplunk.com:587"

// 	if err := smtp.SendMail(addr, auth, "dhanush@geektrust.in", reciever, msg); err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func parseTemplate(templateFileName string, data interface{}) error {
// 	t, err := template.ParseFiles(templateFileName)
// 	if err != nil {
// 		return err
// 	}
// 	buf := new(bytes.Buffer)
// 	if err = t.Execute(buf, data); err != nil {
// 		return err
// 	}
// 	r.body = buf.String()
// 	return nil
// }
