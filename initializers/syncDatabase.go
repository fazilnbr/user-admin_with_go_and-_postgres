package initializer

import (
	"makeconnection.net/sqlandgo/models"
) //user defined packages

func SincDatabase() {
	DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.Admin{})
}
