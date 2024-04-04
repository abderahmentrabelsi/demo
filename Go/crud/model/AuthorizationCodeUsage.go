package model

import (
	"gorm.io/gorm"
	"time"
)

type AuthorizationCodeUsage struct {
	gorm.Model
	Code         string
	Used         bool
	FirstSeenAt  time.Time
	TransientErr bool   // True if a transient error was encountered
	State        string `gorm:"type:varchar(100);default:'Pending'"`
}
