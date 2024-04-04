package controller

import (
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"` // Ensure fields are exported by starting with a capital letter
		Password string `json:"Password"`
	}

	// Use 'err := c.Bind(&body); err != nil' to correctly handle the binding error
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return // Return here to stop function execution after an error
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return // Properly handling the error and stopping the execution
	}

	user := model.User{Email: body.Email, Password: string(hash)}
	result := orm.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return // Properly handling the error and stopping the execution
	}

	// This will be the only response sent to the client if everything above succeeds
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	var user model.User
	orm.DB.Where("email = ?", body.Email).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set JWT token in a cookie
	// Set JWT token in a cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // The cookie is not accessible via JavaScript; helps mitigate XSS attacks
	})

	//log.Println("JWT token successfully stored in a cookie") // Log the success

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in",
		"token":   tokenString, // Return the token to the client
	})
}
