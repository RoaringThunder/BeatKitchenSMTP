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
	"github.com/gorilla/mux"
)

func SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	issuer, ok, err := getTokenIfValid(r)
	params := mux.Vars(r)
	forceSend := params["forceSend"]
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

	user, err := repository.GetUnverifiedUser(issuer, gormDB)
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
	err = service.SendEmail([]string{user.Email}, user.VerificationCode, processedHTML, forceSend)
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
	domain := os.Getenv("COOKIE_DOMAIN")
	cookie, err := r.Cookie(domain)
	if err != nil || cookie == nil {
		utils.HTTPHandleError(w, 400, "Please login and try again")
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SESSION_SECRET")), nil
	})
	if err != nil {
		logging.Log("Failed to parse token: " + err.Error())
		utils.HTTPHandleError(w, 500, "Looks like we're having some issues right now")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		utils.HTTPHandleError(w, 400, "Please login and try again")
		return
	}
	var dbUser models.SalamanderUser
	err = gormDB.Model(&models.SalamanderUser{}).Where("email = ?", user.Email).Find(&dbUser).Error
	if err != nil {
		logging.Log("Bad verify attempt: " + err.Error())
		utils.HTTPHandleError(w, 400, "Bad request")
		return
	}
	if !ok || !token.Valid {
		logging.Log("Couldn't find user: " + user.Email + "Error: " + err.Error())
		utils.HTTPHandleError(w, 500, "Internal server error")
		return
	}
	err = gormDB.Model(&models.VerificationEmailEvent{}).Where("recipient = ? AND  verification_code = ?", user.Email, user.VerificationCode).Updates(map[string]interface{}{"status": "VISITED"}).Error
	if err != nil {
		logging.Log("Failed to update verification email event: " + err.Error())
		utils.HTTPHandleError(w, 500, "Looks like we're having some issues right now")
		return
	}
	fmt.Println(user.Email, "=", claims["iss"], "and!", user.VerificationCode, "=", dbUser.VerificationCode, "and!", dbUser.Status, "!=", "verified")
	if user.Email != claims["iss"].(string) || user.VerificationCode != dbUser.VerificationCode || dbUser.Status == "verified" {
		logging.Log("Bad verify attempt: " + err.Error())
		utils.HTTPHandleError(w, 400, "Bad request")
		return
	}

	fmt.Println(user.Email, user.VerificationCode)

	err = gormDB.Model(&models.SalamanderUser{}).Where("email = ?", user.Email).Updates(map[string]interface{}{"status": "verified"}).Error
	if err != nil {
		logging.Log("Failed to verify user: " + err.Error())
		utils.HTTPHandleError(w, 500, "Looks like we're having some issues right now")
		return
	}

	payload := map[string]interface{}{
		"status":  true,
		"message": "Welcome to Salamander town!",
	}
	utils.HTTPHandleResponse(w, payload)
	return
}

func getTokenIfValid(r *http.Request) (string, bool, error) {
	domain := os.Getenv("COOKIE_DOMAIN")
	cookie, err := r.Cookie(domain)
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
