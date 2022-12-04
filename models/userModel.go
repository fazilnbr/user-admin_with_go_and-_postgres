package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `grom:"unique"`
	Password string
}
