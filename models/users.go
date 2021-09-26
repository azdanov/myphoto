package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrResNotFound is returned when a resource cannot be found in the database.
	ErrResNotFound = errors.New("models: resource was not found")
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: provided ID was invalid")
)

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type UserService struct {
	db *gorm.DB
}

// ByID will find a user by the id.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will find a user by email.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// first will find the first matching record,
// and will place it into dst and return nil.
// If not found returns ErrResNotFound.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResNotFound
	}
	return err
}

// Create will create a new user and fill additional
// fields: ID, CreatedAt and UpdatedAt.
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update will update the user with the provided data.
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will update the user with the provided data.
func (us *UserService) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// DestructiveReset drops the user table and re-creates it.
func (us *UserService) DestructiveReset() error {
	if err := us.db.Migrator().DropTable(&User{}); err != nil {
		panic(err)
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table.
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	return nil
}

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"not null;uniqueIndex"`
	Password string
}
