package entity

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`
	PublishedYear time.Time `json:"published-year" gorm:"type:date;not null"`
	ISBN          string    `json:"isbn" gorm:"type:varchar(255);not null"`
}

type BookReq struct {
	Title         string `json:"title" gorm:"type:varchar(255);not null"`
	PublishedYear string `json:"published-year" gorm:"type:varchar(255);not null"`
	ISBN          string `json:"isbn" gorm:"type:varchar(255);not null"`
}

type BookRead struct {
	gorm.Model
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`
	PublishedYear time.Time `json:"published-year" gorm:"type:date;not null"`
	ISBN          string    `json:"isbn" gorm:"type:varchar(255);not null"`
	Authors       []*Author `json:"authors" gorm:"many2many:mappings;foreignKey:BookReadID;joinForeignKey:BookID;References:ID;joinReferences:AuthorID"`
}
