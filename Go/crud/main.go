package main

import (
	"github.com/cdfmlr/crud/controller"
	"github.com/cdfmlr/crud/middleware"
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Connect to the database and register models
	orm.ConnectDB(orm.DBDriverSqlite, "todolist.db")
	orm.RegisterModel(model.Todo{}, model.Project{}, model.User{}, model.AuthorizationCodeUsage{})

	// Initialize Gin router
	r := gin.Default()

	r.Static("/qrcodes", "./static/qrcodes")

	// Public routes
	publicRoutes := r.Group("/")
	{
		publicRoutes.POST("/signup", controller.SignUp)      // Assuming you have this handler defined
		publicRoutes.POST("/login", controller.LoginHandler) // Updated to use new LoginHandler

		// Adding refresh and logout routes
		publicRoutes.POST("/refresh", controller.RefreshTokenHandler)
		publicRoutes.POST("/logout", controller.LogoutHandler)

		// Setup OAuth2 routes without the AuthMiddleware
		setupOAuth2Routes(publicRoutes)
	}

	// Protected routes with AuthMiddleware
	protectedRoutes := r.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		router.Crud[model.Todo](protectedRoutes, "/todos")
		router.Crud[model.Project](protectedRoutes, "/projects", router.CrudNested[model.Project, model.Todo]("todos"))

		// Add more protected routes here
	}

	// Configure and apply CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the HTTP server on port 8086
	http.ListenAndServe(":8086", handler)
}

// setupOAuth2Routes remains unchanged
func setupOAuth2Routes(r *gin.RouterGroup) {
	r.GET("/auth", controller.AuthHandler)             // Adjust this according to your actual AuthHandler
	r.GET("/callback", controller.AuthCallbackHandler) // Adjust this according to your actual AuthCallbackHandler
}
