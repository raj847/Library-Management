package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(50);not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
}

type UserLoginReg struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
