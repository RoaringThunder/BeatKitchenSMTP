package repository

import (
	"fmt"
	"salamander-smtp/database"
	"salamander-smtp/logging"
	"salamander-smtp/models"
	"strings"

	"gorm.io/gorm"
)

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

func GetUserByEmail(email string, db *gorm.DB) (models.SalamanderUser, error) {
	var user models.SalamanderUser
	err := db.Model(&models.SalamanderUser{}).Where("email = ?", email).Find(&user).Error
	if err != nil {
		logging.Log(fmt.Sprintf("Error fetching user: %s", err))
		return models.SalamanderUser{}, err
	}
	return user, nil
}

func ProcessHTMLTemplate(html string, user models.SalamanderUser) string {
	href := fmt.Sprintf(`<a href="http://localhost:3000/verify/email/%s/v_code/%s">Verify Account</a>`, user.Email, user.VerificationCode)
	processedHTML := strings.ReplaceAll(html, "{{.}}", href)
	return processedHTML
}
