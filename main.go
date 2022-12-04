package main

import (

	// this packages provide by go-internal packages

	// this packages is import from out side - external packages

	"github.com/gin-gonic/gin"
	"makeconnection.net/sqlandgo/controllers"
	initializer "makeconnection.net/sqlandgo/initializers" //user defined packages
	midileware "makeconnection.net/sqlandgo/middleware"
)

// init function work before main function

func init() {

	initializer.LoadEnvVariable()
	initializer.ConnectToDb()
	initializer.SincDatabase()
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", midileware.RequreAuth, controllers.UserHome)

	r.GET("/login", controllers.UserLogin)

	r.POST("/login-submit", controllers.UserAuth)

	r.GET("/register", controllers.UserRegister)

	r.POST("/register-submit", controllers.UserRegisterSubmit)

	r.POST("/signup", controllers.Signup)

	r.POST("/login", controllers.Login)

	r.GET("/validate", midileware.RequreAuth, controllers.Validate)

	r.Run() // listen and serve on 0.0.0.0:8080

}
