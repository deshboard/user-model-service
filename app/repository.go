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
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = grpc.Errorf(codes.NotFound, "user not found")

	// ErrUserAlreadyExists is returned when a user already exists
	ErrUserAlreadyExists = grpc.Errorf(codes.AlreadyExists, "user already exists")
)

// UserRepository implements the gRPC server
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new service object
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

// Create implements the user creation method of the UserRepositoryServer interface
func (s *UserRepository) Create(ctx context.Context, newUser *user.NewUser) (*user.UserCreated, error) {
	var existingID int

	err := s.db.QueryRow("SELECT id FROM users WHERE username = ?", newUser.Username).Scan(&existingID)
	if err != sql.ErrNoRows {
		return nil, ErrUserAlreadyExists
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "%v", err)
	}

	u := &user.User{
		Username:          newUser.Username,
		EncryptedPassword: string(encryptedPassword),
	}

	result, err := s.db.NamedExec(
		"INSERT INTO users (username, encrypted_password) VALUES (:username, :encrypted_password)",
		u,
	)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "%v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "%v", err)
	}

	return &user.UserCreated{Id: id}, nil
}

// Get implements the user lookup method of the UserRepositoryServer interface
func (s *UserRepository) Get(ctx context.Context, getUser *user.GetUser) (*user.User, error) {
	var user user.User
	var err error

	if userID := getUser.GetId(); userID != 0 {
		err = s.db.Get(&user, "SELECT * FROM users WHERE id = ?", userID)
	} else {
		err = s.db.Get(&user, "SELECT * FROM users WHERE username = ?", getUser.GetUsername())
	}

	if err == sql.ErrNoRows {
		err = ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
