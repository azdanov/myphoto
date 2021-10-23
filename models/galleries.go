package models

import "gorm.io/gorm"

// Gallery represents the image resources stored in the database.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}
