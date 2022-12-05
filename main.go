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

	// admin side

	r.GET("/admin", midileware.AdminAuth, controllers.AdminHome)

	r.GET("/admin-login", controllers.AdminLogin)

	r.GET("/admin-logout", controllers.AdminLogout)

	r.POST("/admin-login-submit", controllers.AdminLoginSubmit)

	r.GET("/show-users", midileware.AdminAuth, controllers.AdminShowUser)

	r.GET("/block/:id", midileware.AdminAuth, controllers.AdminUserBlock)

	r.GET("/unblock/:id", midileware.AdminAuth, controllers.AdminUserUnBlock)

	r.GET("/blocked-users", midileware.AdminAuth, controllers.AdminShowBlocedUser)

	r.GET("/adminuserview/:id", midileware.AdminAuth, controllers.AdminUserProfile)

	// user side

	r.GET("/", midileware.RequreAuth, controllers.UserHome)

	r.GET("/login", controllers.UserLogin)

	r.GET("/logout", controllers.UserLogout)

	r.POST("/login-submit", controllers.UserAuth)

	r.POST("/logout", controllers.UserAuth)

	r.GET("/register", controllers.UserRegister)

	r.POST("/register-submit", controllers.UserRegisterSubmit)

	// r.POST("/signup", controllers.Signup)

	// r.POST("/login", controllers.Login)

	// r.GET("/validate", midileware.RequreAuth, controllers.Validate)

	r.Run() // listen and serve on 0.0.0.0:8080

}
