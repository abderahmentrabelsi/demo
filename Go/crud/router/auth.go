package router

import (
	"github.com/cdfmlr/crud/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Add these lines
	r.GET("/auth", controller.AuthHandler)
	r.GET("/callback", controller.AuthCallbackHandler)

	return r
}
