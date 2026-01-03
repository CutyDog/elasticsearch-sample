package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UID string `gorm:"unique;not null"`
}
