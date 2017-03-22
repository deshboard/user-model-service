// +build integration

package app_test

import (
	"testing"

	user "github.com/deshboard/user-model-service/apis/user/v1alpha"
	"github.com/deshboard/user-model-service/app"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
)

func TestDbUserRepository_Create(t *testing.T) {
	service := app.NewService(app.DB)
	ctx := context.Background()

	newUser := &user.NewUser{
		Username: "new",
		Password: "password",
	}

	userCreated, err := service.Create(ctx, newUser)
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	var user user.User

	err = app.DB.Get(&user, "SELECT * FROM users WHERE id = ?", userCreated.GetId())
	if err != nil {
		t.Fatalf("cannot get user, received: %v", err)
	}

	if got, want := user.GetId(), userCreated.GetId(); got != want {
		t.Errorf("user IDs do not match: %d != %d", got, want)
	}

	if got, want := user.GetUsername(), "new"; got != want {
		t.Errorf("username is expected to be '%s', received: %s", want, got)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.GetEncryptedPassword()), []byte("password")); err != nil {
		t.Error("the password check is expected to be successful")
	}
}

func TestDbUserRepository_Create_AlreadyExists(t *testing.T) {
	service := app.NewService(app.DB)
	ctx := context.Background()

	_, err := app.DB.Exec("INSERT INTO users (username, encrypted_password) VALUES ('already_exists', 'encrypted_password')")
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	newUser := &user.NewUser{
		Username: "already_exists",
		Password: "password",
	}

	userCreated, err := service.Create(ctx, newUser)

	if userCreated != nil {
		t.Error("the user should not be created")
	}

	if err != app.ErrUserAlreadyExists {
		t.Error("the user should already exist")
	}
}

func TestDbUserRepository_Get_ById(t *testing.T) {
	service := app.NewService(app.DB)
	ctx := context.Background()

	result, err := app.DB.Exec("INSERT INTO users (username, encrypted_password) VALUES ('get_user_by_id', 'encrypted_password')")
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("cannot get last inserted id, received: %v", err)
	}

	getUser := &user.GetUser{UserKey: &user.GetUser_Id{Id: id}}

	user, err := service.Get(ctx, getUser)
	if err != nil {
		t.Fatalf("cannot get user, received: %v", err)
	}

	if got, want := user.GetId(), id; got != want {
		t.Errorf("user IDs do not match: %d != %d", got, want)
	}

	if got, want := user.GetUsername(), "get_user_by_id"; got != want {
		t.Errorf("username is expected to be '%s', received: %s", want, got)
	}

	if got, want := user.GetEncryptedPassword(), "encrypted_password"; got != want {
		t.Errorf("password is expected to be '%s', received: %s", want, got)
	}
}

func TestDbUserRepository_Get_ByUsername(t *testing.T) {
	service := app.NewService(app.DB)
	ctx := context.Background()

	_, err := app.DB.Exec("INSERT INTO users (username, encrypted_password) VALUES ('get_user_by_username', 'encrypted_password')")
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	getUser := &user.GetUser{UserKey: &user.GetUser_Username{Username: "get_user_by_username"}}

	user, err := service.Get(ctx, getUser)
	if err != nil {
		t.Fatalf("cannot get user, received: %v", err)
	}

	if got, want := user.GetId(), int64(0); got == want {
		t.Error("user ID should not be 0")
	}

	if got, want := user.GetUsername(), "get_user_by_username"; got != want {
		t.Errorf("username is expected to be '%s', received: %s", want, got)
	}

	if got, want := user.GetEncryptedPassword(), "encrypted_password"; got != want {
		t.Errorf("password is expected to be '%s', received: %s", want, got)
	}
}

func TestDbUserRepository_Get_NotFound(t *testing.T) {
	service := app.NewService(app.DB)
	ctx := context.Background()

	getUser := &user.GetUser{UserKey: &user.GetUser_Username{Username: "get_non_existing_user"}}

	user, err := service.Get(ctx, getUser)

	if user != nil {
		t.Error("the user should not be returned")
	}

	if err != app.ErrUserNotFound {
		t.Error("the user should not be found")
	}
}
