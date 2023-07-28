package app

import (
	"salamander-smtp/repository"
	"salamander-smtp/service"
)

func Run() error {
	err, users := repository.GetUnverifiedUsers()
	if err != nil {
		return err
	}
	err, html := repository.GenerateHTML()
	if err != nil {
		return err
	}
	for _, user := range users {
		err = service.SendEmail([]string{user.Email}, html)
		if err != nil {
			return err
		}
	}

	return err
}
