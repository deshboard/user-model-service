package app

import (
	"database/sql"

	user "github.com/deshboard/user-service/apis/iam/user/v1alpha1"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	// ErrUnauthorized is returned when a user cannot be authenticated
	ErrUnauthorized = grpc.Errorf(codes.Unauthenticated, "user cannot be authenticated")
)

// UserDirectory implements the gRPC server
type UserDirectory struct {
	db *sqlx.DB
}

// NewUserDirectory creates a new service object
func NewUserDirectory(db *sqlx.DB) *UserDirectory {
	return &UserDirectory{db}
}

// Authenticate implements the user authentication method of the UserDirectoryServer interface
func (s *UserDirectory) Authenticate(ctx context.Context, credentials *user.Credentials) (*user.AuthenticationResponse, error) {
	var userKey string
	var encryptedPassword string

	err := s.db.QueryRow("SELECT username, encrypted_password FROM users WHERE username = ? OR email = ?", credentials.GetUserKey(), credentials.GetUserKey()).Scan(&userKey, &encryptedPassword)
	if err == sql.ErrNoRows {
		return nil, ErrUnauthorized
	} else if err != nil {
		// TODO: log error
		return nil, ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(credentials.GetPassword()))
	if err != nil {
		return nil, ErrUnauthorized
	}

	return &user.AuthenticationResponse{UserKey: userKey}, nil
}
