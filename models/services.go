package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

func NewServices(psqlInfo string) (*Services, error) {
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

// Close will close the database connection.
func (s *Services) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DestructiveReset will drop all tables and rebuild them.
func (s *Services) DestructiveReset() error {
	if err := s.db.Migrator().DropTable(&User{}, &Gallery{}); err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{})
}
