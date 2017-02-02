package main

import "errors"

// User object returned
type User struct {
	Id                int    `json:"id" db:"id"`
	Username          string `json:"username" db:"username"`
	EncryptedPassword string `json:"encrypted_password" db:"encrypted_password"`
}

// Newly created user representation
type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ErrUserNotFound = errors.New("user: not found")
var ErrUserAlreadyExists = errors.New("user: already exists")

type UserRepository interface {
	// Creates a new user and returns the id of it
	Create(newUser *NewUser) (int, error)

	// Returns a user based on a primary identifier (user ID, username)
	Get(userKey string) (*User, error)
}
