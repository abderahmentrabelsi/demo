package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cdfmlr/crud/config"
	model "github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func AuthHandler(c *gin.Context) {
	url := config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func AuthCallbackHandler(c *gin.Context) {
	ctx := context.Background()
	code := c.Query("code")
	if code == "" {
		log.Println("AuthCallbackHandler error: Code not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
		return
	}

	token, err := config.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("AuthCallbackHandler error: Failed to exchange token: %v\n", err)

		// Specifically handle the "invalid_grant" error more gracefully
		if strings.Contains(err.Error(), "invalid_grant") {
			// Optionally, redirect the user to the login page or show a custom error message
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired authorization code. Please try logging in again."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		}

		return
	}

	user, err := fetchUserInfo(ctx, token)
	if err != nil {
		log.Printf("AuthCallbackHandler error: Failed to fetch user info: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	user, err = ensureUserExists(user, token)
	if err != nil {
		log.Printf("AuthCallbackHandler error: Failed to ensure user exists: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure user exists"})
		return
	}

	sessionToken := uuid.New().String()

	// Add more detailed logging or handling as necessary

	c.JSON(http.StatusOK, gin.H{
		"message":       "User logged in successfully",
		"user":          user,
		"session_token": sessionToken,
	})
}

// fetchUserInfo uses the OAuth2 token to fetch user info from the OAuth provider.
// This needs to be implemented based on your specific OAuth provider's API.
func fetchUserInfo(ctx context.Context, token *oauth2.Token) (*model.User, error) {
	client := config.OAuth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email     string `json:"email"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Picture   string `json:"picture"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &model.User{
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Role:      model.Employee,
		Picture:   userInfo.Picture,
	}, nil
}

func ensureUserExists(tempUser *model.User, token *oauth2.Token) (*model.User, error) {
	var user model.User
	result := orm.DB.Where("email = ?", tempUser.Email).First(&user)

	// Check if the error is because the record was not found, which is an expected scenario
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record not found, so create a new user
			log.Printf("User does not exist, creating a new one: %s\n", tempUser.Email)
			tempUser.AccessToken = token.AccessToken
			tempUser.RefreshToken = token.RefreshToken
			tempUser.TokenExpiry = token.Expiry
			if dbc := orm.DB.Create(&tempUser); dbc.Error != nil {
				log.Printf("ensureUserExists error creating user: %v\n", dbc.Error)
				return nil, dbc.Error
			}
			return tempUser, nil
		} else {
			// An actual error occurred while querying the database
			log.Printf("ensureUserExists error: %v\n", result.Error)
			return nil, result.Error
		}
	}

	// If we reach here, it means the user exists, so update the user's token information
	log.Printf("User exists, updating token info: %s\n", user.Email)
	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.TokenExpiry = token.Expiry
	if dbs := orm.DB.Save(&user); dbs.Error != nil {
		log.Printf("ensureUserExists error updating user: %v\n", dbs.Error)
		return nil, dbs.Error
	}
	return &user, nil
}
