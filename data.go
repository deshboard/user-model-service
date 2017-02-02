package main

import (
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Repository with MySQL database persistence
type DbUserRepository struct {
	db *sqlx.DB
}

// Creates a new DB Repository
func NewDbUserRepository(db *sqlx.DB) *DbUserRepository {
	return &DbUserRepository{db}
}

func (r *DbUserRepository) Create(newUser *NewUser) (int, error) {
	var existingId int
	err := r.db.QueryRow("SELECT id FROM users WHERE username = ?", newUser.Username).Scan(&existingId)
	if err != sql.ErrNoRows {
		return 0, ErrUserAlreadyExists
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &User{
		Username:          newUser.Username,
		EncryptedPassword: string(encryptedPassword),
	}

	result, err := r.db.NamedExec(
		"INSERT INTO users (username, encrypted_password) VALUES (:username, :encrypted_password)",
		user,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	return int(id), err
}

func (r *DbUserRepository) Get(userKey string) (*User, error) {
	var user User
	var err error

	if userId, e := strconv.Atoi(userKey); e == nil {
		err = r.db.Get(&user, "SELECT * FROM users WHERE id = ?", userId)
	} else {
		err = r.db.Get(&user, "SELECT * FROM users WHERE username = ?", userKey)
	}

	if err == sql.ErrNoRows {
		err = ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
