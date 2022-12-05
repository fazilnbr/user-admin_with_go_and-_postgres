package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	initializer "makeconnection.net/sqlandgo/initializers"
	"makeconnection.net/sqlandgo/models"
)

func AdminHome(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminHome.html", gin.H{
		"content": "This is an index page...",
		"message": user,
	})

}
func AdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminSignin.html", gin.H{
		"content": "This is an index page...",
	})

}

func AdminLogout(c *gin.Context) {
	tokenString, err := c.Cookie("adm-auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth", tokenString, -1, "", "", false, true)
	location := url.URL{Path: "/admin"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func AdminLoginSubmit(c *gin.Context) {
	// Get the email and password from req body

	var body struct {
		Username string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// Look up reqested user

	var user models.Admin

	// fmt.Print("\n\n email :", body.Email, "\npassword :", body.Password, "\n\n")

	// It equals to : SELECT * FROM users WHERE email = requested email;
	initializer.DB.First(&user, "Username = ?", body.Username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid *user* name or password ",
		})
		return
	}
	if user.Password != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user name or *password* ",
		})
		return
	}

	// Compare sent password with user password hash

	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "invalid user name or password ",
	// 	})
	// 	return
	// }

	// Create token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create tocken",
		})
		fmt.Println(err)
		return
	}

	// set cookie and Send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("adm-auth", tokenString, 3600*24*30, "", "", false, true)

	// c.JSON(http.StatusOK, gin.H{
	// 	"tocken": tokenString,
	// })
	c.HTML(http.StatusOK, "adminHome.html", gin.H{
		"content": "This is an index page...",
		"tocken":  tokenString,
		"message": user,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"content": "This is an index page...",
	// 	"tocken":  tokenString,
	// 	"message": user,
	// })
}

func AdminShowUser(c *gin.Context) {

	var users []models.User
	// Get all records
	initializer.DB.Where("status = ?", "unblocked").Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminShowuser.html", gin.H{
		"content": "This is an index page...",
		"message": user,
		"users":   users,
	})

	// c.JSON(http.StatusOK, gin.H{
	// 	"content": "This is an index page...",
	// 	"message": user,
	// 	"users":   users,
	// })

}

func AdminUserProfile(c *gin.Context) {

	userid := c.Param("id")
	var user models.User
	initializer.DB.First(&user, userid)
	// SELECT * FROM users WHERE id = 10;

	fmt.Println(user)
	userdata, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminUserProfile.html", gin.H{
		"userdata": user,
		"message":  userdata,
	})

	// c.JSON(http.StatusOK, gin.H{
	// 	"userdata": user,
	// 	"message":  userdata,
	// })
}

func AdminUserBlock(c *gin.Context) {
	userid := c.Param("id")
	var user models.User

	// Update with conditions and model value
	initializer.DB.Model(&user).Where("id = ?", userid).Update("status", "blocked")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

	location := url.URL{Path: "/show-users"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func AdminUserUnBlock(c *gin.Context) {
	param := c.Param("id")
	var user models.User

	// Update with conditions and model value
	initializer.DB.Model(&user).Where("id = ?", param).Update("status", "unblocked")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

	location := url.URL{Path: "/blocked-users"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func AdminShowBlocedUser(c *gin.Context) {
	var users []models.User
	// Get all records
	initializer.DB.Where("status = ?", "blocked").Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminShowBlkduser.html", gin.H{
		"content": "This is an index page...",
		"message": user,
		"users":   users,
	})

}
