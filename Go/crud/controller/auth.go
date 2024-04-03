package controller

import (
	"context"
	"encoding/json"
	"github.com/cdfmlr/crud/config"
	model "github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	url := config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func AuthCallbackHandler(c *gin.Context) {
	ctx := context.Background()
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
		return
	}

	// Exchange the code for a token.
	token, err := config.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Use the token to get user information from the OAuth provider.
	user, err := fetchUserInfo(ctx, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	// Check if the user exists or create a new one.
	user, err = ensureUserExists(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure user exists"})
		return
	}

	// Generate a session token
	sessionToken := uuid.New().String()

	// In a real application, you'd now associate this sessionToken with user data in your session store

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
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &model.User{
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Role:      model.Employee,
	}, nil
}

func ensureUserExists(tempUser *model.User) (*model.User, error) {
	var user model.User
	result := orm.DB.Where("email = ?", tempUser.Email).First(&user)
	if result.RowsAffected == 0 {
		// User does not exist, create a new one
		orm.DB.Create(&tempUser)
		return tempUser, nil
	}
	// Optionally, update user's token information here
	return &user, nil
}
