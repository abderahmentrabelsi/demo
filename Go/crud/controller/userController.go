// Package controller provides the HTTP handlers for the CRUD application.
package controller

import (
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SignUp is a handler function that creates a new user.
// It expects a JSON body with "Email" and "Password" fields.
// If successful, it responds with a 200 status and a success message.
// SignUp is a handler function that creates a new user with TOTP secret.
func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	var existingUser model.User
	result := orm.DB.Where("email = ?", body.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourAppName",
		AccountName: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating TOTP secret"})
		return
	}

	user := model.User{
		Email:      body.Email,
		Password:   string(hash),
		TOTPSecret: totpKey.Secret(),
	}
	result = orm.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Check and create QR code directory if not exists
	qrCodeDir := filepath.Join("static", "qrcodes")
	if _, err := os.Stat(qrCodeDir); os.IsNotExist(err) {
		os.MkdirAll(qrCodeDir, 0755) // Create the directory with necessary permissions
	}

	// Generate QR code for TOTP
	qrCodeData := totpKey.URL()
	qrFilename := filepath.Join(qrCodeDir, body.Email+".png")
	qrCode, err := qrcode.Encode(qrCodeData, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	err = os.WriteFile(qrFilename, qrCode, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save QR code"})
		return
	}

	qrURL := "/qrcodes/" + filepath.Base(qrFilename)
	c.JSON(http.StatusOK, gin.H{
		"message":   "User created successfully",
		"totpURL":   totpKey.URL(),
		"qrCodeURL": qrURL,
	})
}

// LoginHandler is a handler function that authenticates a user.
// It expects a JSON body with "Email" and "Password" fields.
// If successful, it responds with a 200 status and access and refresh tokens.
// LoginHandler is a handler function that authenticates a user and verifies TOTP.
func LoginHandler(c *gin.Context) {
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
		TOTP     string `json:"TOTP"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Capture IP and Device
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

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

	// Verify the TOTP code
	valid := totp.Validate(body.TOTP, user.TOTPSecret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// For testing: Access token expires in 1 minute
	accessToken, err := generateToken(user.Email, 1*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := generateToken(user.Email, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Create a login history record
	history := model.LoginHistory{
		UserID:      user.ID,
		LoginIP:     clientIP,
		LoginDevice: userAgent,
		LoginTime:   time.Now(),
	}

	orm.DB.Create(&history) // Insert the new login history record into the database

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

	newAccessToken, err := generateToken(userEmail, 1*time.Minute)
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
	accessToken := c.Query("accessToken")
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
