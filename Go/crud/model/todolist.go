package main

import (
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/router"
	"github.com/rs/cors"
	"net/http"
)

type Todo struct {
	orm.BasicModel
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Done   bool   `json:"done"`
}

type Project struct {
	orm.BasicModel
	Title string  `json:"title"`
	Todos []*Todo `json:"todos" gorm:"many2many:project_todos"`
}

func main() {
	orm.ConnectDB(orm.DBDriverSqlite, "todolist.db")
	orm.RegisterModel(Todo{}, Project{})

	r := router.NewRouter()
	router.Crud[Todo](r, "/todos")
	router.Crud[Project](r, "/projects",
		router.CrudNested[Project, Todo]("todos"),
	)

	// Setting up CORS with the "github.com/rs/cors" handler.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // This allows access from your React app domain
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"}, // You might want to restrict this in a production environment
	})

	// Wrap the Gin router with the CORS middleware
	handler := c.Handler(r)

	// Serve using the standard net/http server, applying the CORS middleware
	http.ListenAndServe(":8086", handler)

	// If you were using gin-gonic directly for running the server, you would need to adapt this part,
	// as gin-gonic has its own way to handle middleware that could be leveraged instead of using `http.ListenAndServe`.
}
