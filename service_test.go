package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)

	assert.Equal(t, userRepository, service.userRepository)
}

func TestService_CreateUser(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)

	newUser := &NewUser{
		Username: "user",
		Password: "pass",
	}

	userRepository.On("Create", newUser).Return(1, nil)

	b := []byte(`{"username": "user", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.CreateUser(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "/1", w.HeaderMap.Get("Location"))
	assert.Equal(t, 0, w.Body.Len())

	userRepository.AssertExpectations(t)
}

func TestService_CreateUser_InvalidRequest(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)

	b := []byte("invalid_request")

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.CreateUser(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "", w.HeaderMap.Get("Location"))
	assert.Equal(t, 0, w.Body.Len())

	newUser := &NewUser{
		Username: "user",
		Password: "",
	}

	userRepository.AssertNotCalled(t, "Create", newUser)
}

func TestService_CreateUser_AlreadyExists(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)

	newUser := &NewUser{
		Username: "user",
		Password: "pass",
	}

	userRepository.On("Create", newUser).Return(0, ErrUserAlreadyExists)

	b := []byte(`{"username": "user", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.CreateUser(w, req)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, "", w.HeaderMap.Get("Location"))
	assert.Equal(t, 0, w.Body.Len())

	userRepository.AssertExpectations(t)
}

func TestService_CreateUser_Failure(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)

	newUser := &NewUser{
		Username: "user",
		Password: "pass",
	}

	userRepository.On("Create", newUser).Return(0, errors.New("something went wrong"))

	b := []byte(`{"username": "user", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.CreateUser(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "", w.HeaderMap.Get("Location"))
	assert.Equal(t, 0, w.Body.Len())

	userRepository.AssertExpectations(t)
}

func TestService_GetUser(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)
	service.getParams = func(r *http.Request) map[string]string {
		return map[string]string{
			"userKey": "1",
		}
	}

	user := &User{
		Id:                1,
		Username:          "user",
		EncryptedPassword: "encrypted_pass",
	}

	userRepository.On("Get", "1").Return(user, nil)

	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	service.GetUser(w, req)

	returnedUser := new(User)

	err := json.NewDecoder(w.Body).Decode(returnedUser)
	require.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, user, returnedUser)

	userRepository.AssertExpectations(t)
}

func TestService_GetUser_NotFound(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)
	service.getParams = func(r *http.Request) map[string]string {
		return map[string]string{
			"userKey": "1",
		}
	}

	userRepository.On("Get", "1").Return(nil, ErrUserNotFound)

	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	service.GetUser(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, 0, w.Body.Len())

	userRepository.AssertExpectations(t)
}

func TestService_GetUser_Failure(t *testing.T) {
	userRepository := new(UserRepositoryMock)
	service := NewService(userRepository)
	service.getParams = func(r *http.Request) map[string]string {
		return map[string]string{
			"userKey": "1",
		}
	}

	userRepository.On("Get", "1").Return(nil, errors.New("something went wrong"))

	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	service.GetUser(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, 0, w.Body.Len())

	userRepository.AssertExpectations(t)
}
