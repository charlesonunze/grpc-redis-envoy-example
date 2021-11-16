package model

import "gorm.io/gorm"

type (
	// User - the user model
	User struct {
		gorm.Model
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)
