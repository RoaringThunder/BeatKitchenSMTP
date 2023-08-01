package service

import (
	"fmt"
	"os"
	"salamander-smtp/database"
	"salamander-smtp/logging"
	"salamander-smtp/models"
	"strconv"
	"time"

	"github.com/go-mail/mail"
)

func SendEmail(to []string, v_code string, html string, forceSend string) error {
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

	var emailEvent models.VerificationEmailEvent
	err = gormDB.Model(&models.VerificationEmailEvent{}).Where("recipient = ? AND verification_code = ?", to[0], v_code).First(&emailEvent).Order("DESC").Error
	if err != nil && err.Error() == "record not found" {
		if err := d.DialAndSend(m); err != nil {
			err = gormDB.Model(&models.VerificationEmailEvent{}).Create(&models.VerificationEmailEvent{Recipient: to[0], VerificationCode: v_code, Error_msg: err.Error(), Status: "FAILED"}).Error
			if err != nil {
				logging.Log("Failed to update verification email event: " + err.Error())
			}
			return fmt.Errorf("Looks like we're having some issues right now")
		}
		err = gormDB.Model(&models.VerificationEmailEvent{}).Create(&models.VerificationEmailEvent{Recipient: to[0], VerificationCode: v_code}).Error
		if err != nil {
			logging.Log("Failed to create verification email event: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
	} else if err != nil {
		logging.Log("Failed to find verification email event: " + err.Error())
		return fmt.Errorf("Looks like we're having some issues right now")
	} else {
		//if identical email was sent in the last 12 hours and forceSend is false, do nothing
		if (emailEvent.CreatedAt.After(time.Now().Add(-12*time.Hour)) || emailEvent.UpdatedAt.After(time.Now().Add(-12*time.Hour))) && emailEvent.Status != "FAILED" && forceSend == "false" {
			return fmt.Errorf("Email has already been sent")
		}

		// if more than 2 emails were sent in the last hour, dont send
		var sentCount int64
		err = gormDB.Table("smtp.verification_email_event").Select("*").Where("recipient = ? AND verification_code = ? AND created_at >= NOW() - INTERVAL '1 HOUR' AND status != 'FAILED'", to[0], v_code).Count(&sentCount).Error
		if err != nil {
			logging.Log("Failed to count verification email event: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
		if sentCount > 2 {
			return fmt.Errorf("You've sent too many verification emails in the last hour. Please try again later")
		}
		if err := d.DialAndSend(m); err != nil {
			err = gormDB.Model(&models.VerificationEmailEvent{}).Create(&models.VerificationEmailEvent{Recipient: to[0], VerificationCode: v_code, Error_msg: err.Error(), Status: "FAILED"}).Error
			if err != nil {
				logging.Log("Failed to update verification email event: " + err.Error())
			}
			fmt.Println("Failed to send email: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
		err = gormDB.Model(&models.VerificationEmailEvent{}).Create(&models.VerificationEmailEvent{Recipient: to[0], VerificationCode: v_code, Status: "RESEND"}).Error
		if err != nil {
			fmt.Println(err.Error())
			logging.Log("Failed to update verification email event: " + err.Error())
			return fmt.Errorf("Looks like we're having some issues right now")
		}
	}

	return nil
}
