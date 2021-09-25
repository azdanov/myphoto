package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
