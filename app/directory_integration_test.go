// +build integration

package app_test

import (
	"testing"

	"fmt"

	user "github.com/deshboard/user-service/apis/iam/user/v1alpha1"
	"github.com/deshboard/user-service/app"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
)

func TestUserDirectory_Authenticate(t *testing.T) {
	directory := app.NewUserDirectory(app.DB)
	ctx := context.Background()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("cannot create password, received: %v", err)
	}

	_, err = app.DB.Exec(fmt.Sprintf("INSERT INTO users (username, encrypted_password) VALUES ('authenticate', '%s')", encryptedPassword))
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	credentials := &user.Credentials{
		UserKey:  "authenticate",
		Password: "password",
	}

	authResponse, err := directory.Authenticate(ctx, credentials)

	if err != nil {
		t.Fatalf("cannot authenticate user, received: %v", err)
	}

	if got, want := authResponse.GetUserKey(), "authenticate"; got != want {
		t.Errorf("user keys do not match: %s != %s", got, want)
	}
}

func TestUserDirectory_Authenticate_UserNotFound(t *testing.T) {
	directory := app.NewUserDirectory(app.DB)
	ctx := context.Background()

	credentials := &user.Credentials{
		UserKey:  "authenticate_not_found",
		Password: "password",
	}

	authResponse, err := directory.Authenticate(ctx, credentials)

	if authResponse != nil {
		t.Error("the authentication response should not be returned")
	}

	if err != app.ErrUnauthorized {
		t.Error("the user should not be authenticated")
	}
}

func TestUserDirectory_Authenticate_IncorrectPassword(t *testing.T) {
	directory := app.NewUserDirectory(app.DB)
	ctx := context.Background()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("cannot create password, received: %v", err)
	}

	_, err = app.DB.Exec(fmt.Sprintf("INSERT INTO users (username, encrypted_password) VALUES ('authenticate_incorrect_password', '%s')", encryptedPassword))
	if err != nil {
		t.Fatalf("cannot create user, received: %v", err)
	}

	credentials := &user.Credentials{
		UserKey:  "authenticate_incorrect_password",
		Password: "incorrect_password",
	}

	authResponse, err := directory.Authenticate(ctx, credentials)

	if authResponse != nil {
		t.Error("the authentication response should not be returned")
	}

	if err != app.ErrUnauthorized {
		t.Error("the user should not be authenticated")
	}
}
