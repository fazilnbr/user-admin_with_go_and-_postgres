package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Dob      string
	Gender   string
	Email    string `grom:"unique" gorm:"primaryKey"`
	Password string
	Status   string
}

type Admin struct {
	gorm.Model
	Username string `grom:"unique"`
	Password string
}
