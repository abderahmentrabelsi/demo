// Package controller provides the HTTP handlers for the CRUD application.
package controller

import (
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// SignUp is a handler function that creates a new user.
// It expects a JSON body with "Email" and "Password" fields.
// If successful, it responds with a 200 status and a success message.
func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	// Bind the incoming JSON to body struct
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// Check if user already exists
	var existingUser model.User
	result := orm.DB.Where("email = ?", body.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user if not existing
	user := model.User{Email: body.Email, Password: string(hash)}
	result = orm.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// LoginHandler is a handler function that authenticates a user.
// It expects a JSON body with "Email" and "Password" fields.
// If successful, it responds with a 200 status and access and refresh tokens.
func LoginHandler(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user model.User
	result := orm.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// For testing: Access token expires in 1 minute
	accessToken, err := generateToken(user.Email, 1*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := generateToken(user.Email, 24*time.Hour) // Refresh token with a standard expiration
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	store.SetRefreshToken(refreshToken, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshTokenHandler is a handler function that refreshes a user's access token.
// It expects a form data with "refresh_token" field.
// If successful, it responds with a 200 status and a new access token.
func RefreshTokenHandler(c *gin.Context) {
	refreshToken := c.PostForm("refresh_token")
	userEmail, exists := store.GetEmailByRefreshToken(refreshToken)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, err := generateToken(userEmail, 1*time.Minute) // Ensure short-lived for testing
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// LogoutHandler is a handler function that logs out a user.
// It expects query parameters "accessToken" and "refreshToken".
// If successful, it responds with a 200 status and a success message.
func LogoutHandler(c *gin.Context) {
	accessToken := c.Query("accessToken") // Changed to using query parameter for flexibility
	refreshToken := c.Query("refreshToken")

	store.RevokeToken(accessToken)
	store.RemoveRefreshToken(refreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// generateToken is a helper function that generates a JWT token.
// It takes an email and a duration for the token's expiration.
// It returns the token as a string and any error encountered.
func generateToken(email string, duration time.Duration) (string, error) {
	exp := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   exp.Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

//TODO if con from other browser discon
//TODO actual browser (ip) and device
