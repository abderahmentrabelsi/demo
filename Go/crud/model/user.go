package models

import (
	"github.com/cdfmlr/crud/orm"
	"time"
)

type Role string

const (
	Admin    Role = "Admin"
	Manager  Role = "Manager"
	Employee Role = "Employee"
)

type User struct {
	orm.BasicModel
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Role         Role      `json:"role"`
	WorkHours    int       `json:"work_hours"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
}
