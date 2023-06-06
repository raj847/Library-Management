package entity

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:varchar(255);not null"`
	Country string `json:"country" gorm:"type:varchar(255);not null"`
}

type AuthorReg struct {
	Name    string `json:"name" gorm:"type:varchar(255);not null"`
	Country string `json:"country" gorm:"type:varchar(255);not null"`
}

type AuthorRead struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:varchar(255);not null"`
	Country string `json:"country" gorm:"type:varchar(255);not null"`
	Books   []Book `json:"books" gorm:"many2many:mappings;"`
}
