package models

import (
	"time"

	"gorm.io/gorm"
)

type SalamanderUser struct {
	gorm.Model
	Password         string    `gorm:"not null" json:"password"`
	Email            string    `gorm:"unique;not-null;size:320" json:"email"`
	SpotifyID        string    `gorm:"unique" json:"spotify_id"`
	Status           string    `gorm:"default:'notVerified'" json:"status"`
	VerificationCode string    `gorm:"default:''" json:"verification_code"`
	VCodeStamp       time.Time `json:"vcode_stamp"`
}

func (SalamanderUser) TableName() string {
	return "sd_01.smdr_user"
}
