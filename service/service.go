package service

import (
	"fmt"
	"os"
	"salamander-smtp/database"
	"salamander-smtp/logging"
	"salamander-smtp/models"
	"strconv"

	"github.com/go-mail/mail"
)

func SendEmail(to []string, v_code string, html string) error {
	gormDB := database.FetchDB()
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
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	var emailEvent models.VerificationEmailEvent
	err = gormDB.Model(&models.VerificationEmailEvent{}).Where("recipient = ? AND verification_code = ?", to[0], v_code).First(&emailEvent).Error

	if err != nil && err.Error() == "record not found" {
		err = gormDB.Model(&models.VerificationEmailEvent{}).Create(&models.VerificationEmailEvent{Recipient: to[0], VerificationCode: v_code}).Error
		if err != nil {
			logging.Log("Failed to create verification email event: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
	} else if err != nil {
		logging.Log("Failed to find verification email event: " + err.Error())
		return fmt.Errorf("Looks like we're having some issues right now")
	} else {
		err = gormDB.Model(&models.VerificationEmailEvent{}).Where("recipient = ? AND verification_code = ?", to[0], v_code).Updates(map[string]interface{}{"status": "RESENT"}).Error
		if err != nil {
			logging.Log("Failed to update verification email event: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
	}

	return nil
}
