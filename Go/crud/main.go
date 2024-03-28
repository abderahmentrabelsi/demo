package main

import (
	"github.com/cdfmlr/crud/model"
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/router"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Connect to the database and register models
	orm.ConnectDB(orm.DBDriverSqlite, "todolist.db")
	orm.RegisterModel(models.Todo{}, models.Project{}, models.User{})

	// Setup the HTTP router and routes
	r := router.NewRouter()

	// Configure CRUD operations for Todo, Project, and User models
	router.Crud[models.Todo](r, "/todos")
	router.Crud[models.Project](r, "/projects", router.CrudNested[models.Project, models.Todo]("todos"))
	router.Crud[models.User](r, "/users")

	// Configure and apply CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the HTTP server
	http.ListenAndServe(":8086", handler)
}
