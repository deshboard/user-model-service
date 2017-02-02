package main

import "github.com/stretchr/testify/mock"

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) Create(newUser *NewUser) (int, error) {
	args := r.Called(newUser)

	return args.Int(0), args.Error(1)
}

func (r *UserRepositoryMock) Get(userKey string) (*User, error) {
	args := r.Called(userKey)

	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}

	return user.(*User), args.Error(1)
}
