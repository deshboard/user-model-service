package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Returns parameters from the request
// (decouples the service from the router implementation)
type ParamFetcher func(r *http.Request) map[string]string

type Service struct {
	userRepository UserRepository
	getParams      ParamFetcher
}

// Creates a new service object
func NewService(userRepository UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
		getParams: func(r *http.Request) map[string]string {
			return make(map[string]string)
		},
	}
}

// Handles creating a new user
func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(NewUser)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		// TODO: add logging
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// TODO: add validation
	id, err := s.userRepository.Create(newUser)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			w.WriteHeader(http.StatusConflict)

			return
		default:
			// TODO: add logging
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	w.Header().Set("Location", fmt.Sprintf("/%d", id))
	w.WriteHeader(http.StatusCreated)
}

// Handles retrieving a user based on it's primary identifier
func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	params := s.getParams(r)

	user, err := s.userRepository.Get(params["userKey"])
	if err != nil {
		switch err {
		case ErrUserNotFound:
			w.WriteHeader(http.StatusNotFound)

			return
		default:
			// TODO: add logging
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	json.NewEncoder(w).Encode(user)

	w.WriteHeader(http.StatusOK)
}
