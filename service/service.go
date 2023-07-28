package service

import (
	"os"
	"strconv"

	"github.com/go-mail/mail"
)

func SendEmail(to []string, html string) error {
	username := os.Getenv("SMTP_USERNAME")
	secret := os.Getenv("SMTP_SECRET")
	host := os.Getenv("SMTP_HOST")

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}
	from := os.Getenv("SMTP_FROM")

	m := mail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Verify Email!")
	m.SetBody("text/html", html)

	d := mail.NewDialer(host, port, username, secret)

	// Send the email to recipients.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
