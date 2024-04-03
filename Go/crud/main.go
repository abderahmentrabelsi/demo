package main

import (
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/router"
	"github.com/gin-gonic/gin"
	//model
	"github.com/cdfmlr/crud/controller"
	model "github.com/cdfmlr/crud/model"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Connect to the database and register models
	orm.ConnectDB(orm.DBDriverSqlite, "todolist.db")
	orm.RegisterModel(model.Todo{}, model.Project{}, model.User{})

	// Initialize Gin router
	r := gin.Default()
	router.Crud[model.Todo](r, "/todos")
	router.Crud[model.Project](r, "/projects", router.CrudNested[model.Project, model.Todo]("todos"))
	router.Crud[model.User](r, "/users")

	// Setup OAuth2 routes

	setupOAuth2Routes(r)

	// Existing CRUD operations setup...
	// Note: Replace "router.NewRouter()" with "gin.Default()" based on your provided code.
	// For example, if you have specific configurations in NewRouter(), make sure they are applied here.

	// Configure and apply CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Adjust according to your frontend setup
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Include Authorization for OAuth2 tokens
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the HTTP server on port 8086
	http.ListenAndServe(":8086", handler)
}

// setupOAuth2Routes configures routes related to OAuth2 authentication.
func setupOAuth2Routes(r *gin.Engine) {
	r.GET("/auth", controller.AuthHandler)
	r.GET("/callback", controller.AuthCallbackHandler)
}
