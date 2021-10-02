package models

import (
	"errors"
	"myphoto/hash"
	"myphoto/rand"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrResNotFound is returned when a resource cannot be found in the database.
	ErrResNotFound = errors.New("models: resource was not found")
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: provided ID was invalid")
	// ErrInvalidPassword is returned when an invalid password is used for login.
	ErrInvalidPassword = errors.New("models: provided password was invalid")
)

const hmacSecretKey = "secret-hmac-key"

func NewUserService(db *gorm.DB) *UserService {
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{db: db, hmac: hmac}
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// ByID will find a user by the id.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will find a user by an email.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	err := first(us.db.Where("email = ?", email), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByRemember will find a user by a rememberToken.
func (us *UserService) ByRemember(rememberToken string) (*User, error) {
	var user User
	err := first(us.db.Where("remember_hash = ?", us.hmac.Hash(rememberToken)), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
	p, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = p
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}

// hashPassword will hash a given password with bcrypt and a default cost.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Authenticate will authenticate a user and returns a User on success
// or return an error on failure.
func (us *UserService) Authenticate(email string, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = checkPasswordHash(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// checkPasswordHash compares a password and a hash and
// either returns nil on success or an error on failure.
func checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

// Update will update the user with the provided data.
func (us *UserService) Update(user *User) error {
	p, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = p
	user.Password = ""
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
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
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}
