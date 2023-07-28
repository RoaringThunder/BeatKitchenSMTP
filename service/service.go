package service

import (
	"net/smtp"
	"os"
)

func SendEmail(to []string, html string) error {
	username := os.Getenv("SMTP_USERNAME")
	secret := os.Getenv("SMTP_SECRET")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	addr := host + ":" + port

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Verify email" + "!\n" + "To: " + to[0] + "\n"
	msg := []byte(subject + mime + "\n" + html)

	auth := smtp.CRAMMD5Auth(username, secret)

	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
