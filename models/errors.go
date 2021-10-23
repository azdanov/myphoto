package models

import (
	"strings"
)

const (
	// ErrResourceNotFound is returned when a resource cannot be found in the database.
	ErrResourceNotFound publicError = "resource was not found"

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID privateError = "ID is invalid"

	// ErrInvalidPassword is returned when an invalid password is used for login.
	ErrInvalidPassword publicError = "password is invalid"

	// ErrRequiredEmail is returned when an empty email is provided.
	ErrRequiredEmail publicError = "email address is required"

	// ErrInvalidEmail is returned when an invalid format of email is provided.
	ErrInvalidEmail publicError = "email address is invalid"

	// ErrUnavailableEmail is returned when an email address is already taken.
	ErrUnavailableEmail publicError = "email address is already taken"

	// ErrRequiredPassword is returned when an empty password is provided
	ErrRequiredPassword publicError = "password is required"

	// ErrShortPassword is returned when a passwords' length is too short
	ErrShortPassword publicError = "password length must be at least 8 characters"

	// ErrShortRemember is returned when a remember-tokens' length is too short
	ErrShortRemember privateError = "remember token length must be at least 32 bytes"

	// ErrRequiredRemember is returned when an empty remember token is provided
	ErrRequiredRemember privateError = "remember token is required"
)

type publicError string

func (e publicError) Error() string {
	return string(e)
}

func (e publicError) Public() string {
	split := strings.Split(string(e), " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
