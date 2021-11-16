package model

import "gorm.io/gorm"

type (
	// Account - the user balance model
	Account struct {
		gorm.Model
		ID      string `json:"id"`
		UserID  string `json:"user_id"`
		Balance int    `json:"balance"`
	}
)
