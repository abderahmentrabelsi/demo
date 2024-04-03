package controller

import (
	"github.com/cdfmlr/crud/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	url := config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func AuthCallbackHandler(c *gin.Context) {
	// Implement the logic to handle the callback from the OAuth provider,
	// exchange the code for tokens, and handle user creation or updates.
}
