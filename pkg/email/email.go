package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
}

// NewEmailService khởi tạo dịch vụ gửi email
func NewEmailService(smtpHost string, smtpPort int, username, password string) *EmailService {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

// SendEmail gửi email với các thông tin cơ bản
func (e *EmailService) SendEmail(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.username, e.password)

	// Gửi email
	if err := d.DialAndSend(msg); err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
