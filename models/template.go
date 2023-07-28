package models

import "gorm.io/gorm"

type VerifyUserTemplate struct {
	gorm.Model
	HTML string `json:"html" gorm:"column:html"`
}

func (t *VerifyUserTemplate) TableName() string {
	return "sd_01.verify_user_templates"
}
