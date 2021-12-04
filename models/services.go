package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

type ServicesConfig func(*Services) error

func WithGorm(dsn string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser(hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func NewServices(configs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, config := range configs {
		if err := config(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
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
