package models

import (
	"errors"
	"myphoto/hash"
	"myphoto/rand"
	"strings"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	hmacSecretKey     = "secret-hmac-key"
	minPasswordLength = 8
)

// User represents the user model stored in the database.
// Used for user accounts, storing both an email and a
// password so users can log in and gain access to content.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}

// UserDB is used to interact with the users' database.
//
// For the majority of single user queries:
// If the user is found then error is nil.
// If the user is not found error is ErrResourceNotFound.
// If there is another error then it is returned.
//
// For single user queries any error apart from ErrResourceNotFound
// should result in a 500 status code.
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(rememberToken string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// UserService is a set of methods used to manipulate
// and work with the user model.
type UserService interface {
	UserDB
	// Authenticate will verify the provided email and
	// password are correct. If they are correct, the
	// User corresponding to that email is returned.
	// Otherwise, either ErrResourceNotFound, ErrInvalidPassword or
	// another error.
	Authenticate(email, password string) (*User, error)
}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db: db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{UserDB: ug, hmac: hmac}
	return &userService{UserDB: uv}
}

// Confirm that userService implements UserDB interface.
var _ UserDB = &userService{}

type userService struct {
	UserDB
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}
	return user, nil
}

// Confirm that userValidator implements UserDB interface.
var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	var user User
	user.Email = email
	err := runUserValFuncs(&user, uv.normalizeEmail, uv.requireEmail, uv.validateEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	var user User
	user.Remember = token
	err := runUserValFuncs(&user, uv.hashRemember)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.validateEmail,
		uv.availableEmail,
		uv.requiredPassword,
		uv.validatePassword,
		uv.hashPassword,
		uv.requiredPasswordHash,
		uv.ensureRemember,
		uv.validateRemember,
		uv.hashRemember,
		uv.requiredRememberHash,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.validateEmail,
		uv.availableEmail,
		uv.validatePassword,
		uv.hashPassword,
		uv.requiredPasswordHash,
		uv.validateRemember,
		uv.hashRemember,
		uv.requiredRememberHash,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(user.ID)
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) hashPassword(u *User) error {
	if u.Password != "" {
		pBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(pBytes)
		u.Password = ""
	}
	return nil
}

func (uv *userValidator) requiredPassword(u *User) error {
	if u.Password == "" {
		return ErrRequiredPassword
	}
	return nil
}

func (uv *userValidator) requiredPasswordHash(u *User) error {
	if u.PasswordHash == "" {
		return ErrRequiredPassword
	}
	return nil
}

func (uv *userValidator) validatePassword(u *User) error {
	if u.Password != "" {
		if len([]rune(u.Password)) < minPasswordLength {
			return ErrShortPassword
		}
	}
	return nil
}

func (uv *userValidator) ensureRemember(u *User) error {
	if u.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		u.Remember = token
	}
	return nil
}

func (uv *userValidator) validateRemember(u *User) error {
	if u.Remember != "" {
		n, err := rand.NBytes(u.Remember)
		if err != nil {
			return err
		}
		if n < rand.RememberTokenBytes {
			return ErrShortRemember
		}
	}
	return nil
}

func (uv *userValidator) hashRemember(u *User) error {
	if u.Remember != "" {
		u.RememberHash = uv.hmac.Hash(u.Remember)
	}
	return nil
}

func (uv *userValidator) requiredRememberHash(u *User) error {
	if u.RememberHash == "" {
		return ErrRequiredRemember
	}
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return func(u *User) error {
		if u.ID <= n {
			return ErrInvalidID
		}
		return nil
	}
}

func (uv *userValidator) requireEmail(u *User) error {
	if u.Email == "" {
		return ErrRequiredEmail
	}
	return nil
}

func (uv *userValidator) normalizeEmail(u *User) error {
	u.Email = strings.TrimSpace(u.Email)
	u.Email = strings.ToLower(u.Email)
	return nil
}

func (uv *userValidator) validateEmail(u *User) error {
	err := checkmail.ValidateFormat(u.Email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func (uv *userValidator) availableEmail(u *User) error {
	// Warning not to use this validator inside ByEmail to avoid cyclic call
	existingUser, err := uv.ByEmail(u.Email)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			// Email is available
			return nil
		}
		return err
	}
	if existingUser.ID != u.ID {
		// Users are different; Not an update of email by same user
		return ErrUnavailableEmail
	}
	return nil
}

// Confirm that userGorm implements UserDB interface.
var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	err := first(ug.db.Where("email = ?", email), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResourceNotFound
	}
	return err
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}
