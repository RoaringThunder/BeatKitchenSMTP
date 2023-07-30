package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"salamander-smtp/database"
	"salamander-smtp/logging"
	"salamander-smtp/models"
	"salamander-smtp/repository"
	"salamander-smtp/service"
	utils "salamander-smtp/utils/responseUtils"

	"github.com/golang-jwt/jwt/v5"
)

func SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	issuer, ok, err := getTokenIfValid(r)
	if err != nil {
		utils.HTTPHandleError(w, 400, err.Error())
		return
	}
	if !ok {
		utils.HTTPHandleError(w, 400, "Invalid cookie")
		return
	}

	gormDB := database.FetchDB()
	if err != nil {
		logging.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := repository.GetUserByEmail(issuer, gormDB)
	if err != nil {
		utils.HTTPHandleError(w, 500, err.Error())
		return
	}
	err, html := repository.GenerateHTML()
	if err != nil {
		utils.HTTPHandleError(w, 500, err.Error())
		return
	}
	processedHTML := repository.ProcessHTMLTemplate(html, user)
	if err != nil {
		utils.HTTPHandleError(w, 500, err.Error())
		return
	}
	err = service.SendEmail([]string{user.Email}, processedHTML)
	if err != nil {
		utils.HTTPHandleError(w, 500, err.Error())
		return
	}
	payload := map[string]interface{}{
		"status":  true,
		"message": fmt.Sprintf("Sent email successfully to: %s", user.Email),
	}
	utils.HTTPHandleResponse(w, payload)
	return
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	gormDB := database.FetchDB()
	var user models.SalamanderUser
	err := json.NewDecoder(r.Body).Decode(&user)
	cookie, err := r.Cookie("salamander.api")
	if err != nil {
		utils.HTTPHandleError(w, 400, "No cookie")
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SESSION_SECRET")), nil
	})
	if err != nil {
		utils.HTTPHandleError(w, 400, "Error reading cookie")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		utils.HTTPHandleError(w, 400, "Invalid cookie")
		return
	}
	var dbUser models.SalamanderUser
	err = gormDB.Model(&models.SalamanderUser{}).Where("email = ?", user.Email).Find(&dbUser).Error
	if !ok || !token.Valid {
		logging.Log("Couldn't find user: " + user.Email + "Error: " + err.Error())
		utils.HTTPHandleError(w, 500, "Internal server error")
		return
	}
	if user.Email != claims["iss"] || user.VerificationCode != dbUser.VerificationCode || dbUser.Status == "verified" {
		fmt.Println(claims["iss"], user.Email)
		fmt.Println(user.VerificationCode, dbUser.VerificationCode)
		utils.HTTPHandleError(w, 400, "Bad request")
		return
	}

	gormDB.Model(&models.SalamanderUser{}).Where("email = ?", user.Email).Updates(map[string]interface{}{"status": "verified"})
	payload := map[string]interface{}{
		"status":  true,
		"message": "Welcome to Salamander town!",
	}
	utils.HTTPHandleResponse(w, payload)
	return
}

func getTokenIfValid(r *http.Request) (string, bool, error) {
	cookie, err := r.Cookie("salamander.api")
	if err != nil {
		return "", false, err
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SESSION_SECRET")), nil
	})
	if err != nil {
		return "", false, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", false, err
	}
	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", false, fmt.Errorf("unable to extract issuer claim")
	}
	return issuer, true, nil
}
