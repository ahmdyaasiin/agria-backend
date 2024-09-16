package verification

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, message string) error {
	SmtpHost := os.Getenv("SMTP_HOST")
	SmtpPort := os.Getenv("SMTP_PORT")
	SmtpUsername := os.Getenv("SMTP_USERNAME")
	SmtpPassword := os.Getenv("SMTP_PASSWORD")

	addr := fmt.Sprintf("%s:%s", SmtpHost, SmtpPort)
	msg := fmt.Sprintf("From: Test Email <%s>\nTo: %s\nSubject:%s\n\n%s", SmtpUsername, to, subject, message)
	err := smtp.SendMail(addr,
		smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpHost),
		SmtpUsername, []string{to}, []byte(msg))

	return err
}
