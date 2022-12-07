package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Dob      string
	Gender   string
	Mobile   string `gorm:"index:idx_name,unique"`
	Email    string `gorm:"index:idx_name,unique"`
	Password string
	Status   string
}

type Admin struct {
	gorm.Model
	Username string `grom:"unique"`
	Password string
}
