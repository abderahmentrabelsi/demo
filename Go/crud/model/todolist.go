package models

import (
	"github.com/cdfmlr/crud/orm"
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

//new model for user
