package main

import (
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewDbUserRepository_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	assert.Equal(t, db, userRepository.db)
}

func TestDbUserRepository_Create_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	newUser := &NewUser{
		Username: "new",
		Password: "password",
	}

	id, err := userRepository.Create(newUser)

	assert.NoError(t, err)

	var user User

	err = db.Get(&user, "SELECT * FROM users WHERE id = ?", id)

	assert.NoError(t, err)

	assert.Equal(t, id, user.Id)
	assert.Equal(t, "new", user.Username)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte("password")))
}

func TestDbUserRepository_Create_AlreadyExists_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	_ = db.MustExec("INSERT INTO users (username, encrypted_password) VALUES ('already_exists', 'encrypted_password')")

	newUser := &NewUser{
		Username: "already_exists",
		Password: "password",
	}

	id, err := userRepository.Create(newUser)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, ErrUserAlreadyExists, err)
}

func TestDbUserRepository_Get_ById_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	result := db.MustExec("INSERT INTO users (username, encrypted_password) VALUES ('get_user_by_id', 'encrypted_password')")

	id, err := result.LastInsertId()
	require.NoError(t, err)

	userKey := strconv.FormatInt(id, 10)

	user, err := userRepository.Get(userKey)

	assert.NoError(t, err)

	assert.Equal(t, int(id), user.Id)
	assert.Equal(t, "get_user_by_id", user.Username)
	assert.Equal(t, "encrypted_password", user.EncryptedPassword)
}

func TestDbUserRepository_Get_ByUsername_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	db.MustExec("INSERT INTO users (username, encrypted_password) VALUES ('get_user_by_username', 'encrypted_password')")

	user, err := userRepository.Get("get_user_by_username")

	assert.NoError(t, err)

	assert.NotEqual(t, 0, user.Id)
	assert.Equal(t, "get_user_by_username", user.Username)
	assert.Equal(t, "encrypted_password", user.EncryptedPassword)
}

func TestDbUserRepository_Get_NotFound_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	userRepository := NewDbUserRepository(db)

	user, err := userRepository.Get("get_non_existing_user")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrUserNotFound, err)
}
