package models

import (
	"errors"

	"gorm.io/gorm"
)

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResourceNotFound
	}
	return err
}
