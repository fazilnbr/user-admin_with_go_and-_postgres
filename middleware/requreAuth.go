package midileware

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

func RequreAuth(c *gin.Context) {

	defer func() {

		if err := recover(); err != nil {
			c.HTML(http.StatusOK, "home.html", gin.H{})
			fmt.Printf("\n\nRecovered from panic. %s\n\n", err)
		}
	}()
	// Get the cookie of the request

	tokenString, err := c.Cookie("auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/validation

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token != nil {

		// check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user with token subject

		var user models.User

		initializer.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			location := url.URL{Path: "/"}
			c.Redirect(http.StatusFound, location.RequestURI())

			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if user.Status == "blocked" {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"error": "Admin Blocked You",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if user.Status == "newuser" {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"error": "You Are Not Activated Yet",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// atach to request

		c.Set("user", user)

		// continue

		c.Next()

		// fmt.Println(claims["exp"], claims["sub"])
	} else {
		c.HTML(http.StatusOK, "home.html", gin.H{})

		//c.AbortWithStatus(http.StatusUnauthorized)
	}

}
