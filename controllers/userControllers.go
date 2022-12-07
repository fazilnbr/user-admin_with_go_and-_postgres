package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	initializer "makeconnection.net/sqlandgo/initializers"
	"makeconnection.net/sqlandgo/models"
)

func Signup(c *gin.Context) {

	// get the email and password of req body

	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// create the user

	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializer.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// respond

	c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	// Get the email and password from req body

	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// Look up reqested user

	var user models.User

	// fmt.Print("\n\n email :", body.Email, "\npassword :", body.Password, "\n\n")

	// It equals to : SELECT * FROM users WHERE email = requested email;
	initializer.DB.First(&user, "Email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user name or password ",
		})
		return
	}

	// Compare sent password with user password hash

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user name or password ",
		})
		return
	}

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
	c.SetCookie("auth", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"tocken": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"meassage": user,
	})
}

func UserHome(c *gin.Context) {
	user, _ := c.Get("user")
	fmt.Println(user)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"message": user,
	})
}

func UserLogin(c *gin.Context) {

	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	c.HTML(http.StatusOK, "signin.html", gin.H{
		"content": "This is an index page...",
	})
}

func UserLogout(c *gin.Context) {

	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	tokenString, err := c.Cookie("auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth", tokenString, -1, "", "", false, true)
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func UserAuth(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	// Get the email and password from req body

	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if c.Bind(&body) != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"error": "Invalid Inputs Please Check Inputs",
		})
		return
	}

	// Look up reqested user

	var user models.User

	// fmt.Print("\n\n email :", body.Email, "\npassword :", body.Password, "\n\n")

	// It equals to : SELECT * FROM users WHERE email = requested email;
	initializer.DB.First(&user, "Email = ?", body.Email)

	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"error": "invalid user name or password ",
		})
		return
	}
	if user.Status == "blocked" {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"error": "you are blocked by admin",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Status == "newuser" {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"error": "You Are Not Activated Yet",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Compare sent password with user password hash

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user name or password ",
		})
		return
	}

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
	c.SetCookie("auth", tokenString, 3600*24*30, "", "", false, true)

	// c.JSON(http.StatusOK, gin.H{
	// 	"tocken": tokenString,
	// })
	c.HTML(http.StatusOK, "home.html", gin.H{
		"content": "This is an index page...",
		"tocken":  tokenString,
		"message": user,
	})
}

func UserRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{
		"content": "This is an index page...",
	})
}

func UserRegisterSubmit(c *gin.Context) {
	// get the email and password of req body

	var body struct {
		Name     string `json:"name" binding:"required"`
		Dob      string
		Gender   string
		Mobile   string `json:"mobile" binding:"required,len=10"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Status   string
	}
	if err := c.Bind(&body); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Invalid Inputs Please Check Inputs",
		})
		return
	}

	// hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// create the user

	user := models.User{Name: body.Name, Dob: body.Dob, Gender: body.Gender, Mobile: body.Mobile, Email: body.Email, Password: string(hash), Status: "newuser"}

	result := initializer.DB.Create(&user) // pass pointer of data to Create

	// fmt.Print("\n\nDB ERROR :", result, "\\n\n")

	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "failed to create user Mobile Number or Email alrady exist Check inputs   ",
		})
		return
	}

	// respond
	c.HTML(http.StatusOK, "signin.html", gin.H{
		"content": "This is an index page...",
	})
}

func UserEditProfile(c *gin.Context) {

	user, _ := c.Get("user")
	fmt.Println("edit user  :", user)
	c.HTML(http.StatusOK, "signup.html", gin.H{
		"content": "This is an index page...",
		"message": user,
	})
}

func UserEditProfileSubmit(c *gin.Context) {
	//get the user
	userid := c.Param("id")
	// userdata, _ := c.Get("user")
	// get the email and password of req body
	var body struct {
		Name   string `json:"name" binding:"required"`
		Dob    string
		Gender string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Inputs Please Check Inputs",
		})
		return
	}
	// create the user
	var user models.User
	// Update with conditions and model value
	initializer.DB.Model(&user).Where("id = ?", userid).Updates(map[string]interface{}{"Name": body.Name, "Dob": body.Dob, "Gender": body.Gender})
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())

}
