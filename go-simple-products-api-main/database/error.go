package database

import "errors"

var (
	// ErrIDNotFound is returned when an id was not found
	ErrIDNotFound = errors.New("provided id was not found")
	// is return when a invalid id
	ErrInvalidID = errors.New("provided id was invalid")
)
