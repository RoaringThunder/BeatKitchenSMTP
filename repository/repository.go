package repository

import (
	"fmt"
	"salamander-smtp/database"
	"salamander-smtp/logging"
	"salamander-smtp/models"
)

func GetUnverifiedUsers() (error, []models.SalamanderUser) {
	var users []models.SalamanderUser
	db := database.FetchDB()

	err := db.Model(&users).Where("status = ?", "notVerified").Find(&users).Error
	if err != nil {
		logging.Log(fmt.Sprintf("Error fetching unverified users: %s", err))
		return err, []models.SalamanderUser{}
	}
	return err, users
}

func GenerateHTML() (error, string) {
	db := database.FetchDB()
	var template models.VerifyUserTemplate
	err := db.Model(&models.VerifyUserTemplate{}).Where("id = ?", 2).Find(&template).Error
	if err != nil {
		logging.Log(fmt.Sprintf("Error fetching template: %s", err))
		return err, ""
	}

	return nil, template.HTML
}
