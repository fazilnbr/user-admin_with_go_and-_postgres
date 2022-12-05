package initializer

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {

	// Provide credencials fot database
	const (
		host     = "localhost"
		port     = 5432
		user     = "fasil"
		database = "useradminauth"
		password = "0000"
	)

	// dns-data sorce name

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

	// Establish connections to the database and returning if any error

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("faild to connect ")

	}
	fmt.Println("connected successfully")

}
