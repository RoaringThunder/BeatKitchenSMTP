package models

import "time"

type VerificationEmailEvent struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Recipient        string    `json:"recipient"`
	Status           string    `json:"status" gorm:"default:'SENT'"`
	VerificationCode string    `json:"verification_code"`
	Error_msg        string    `json:"error_msg"`
}

func (s *VerificationEmailEvent) TableName() string {
	return "smtp.verification_email_event"
}
