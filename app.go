package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db      *sqlx.DB
	service *Service
}

// Creates a new app
func NewApp() (*App, error) {
	db, err := sqlx.Open(
		RequireEnv("DATABASE_TYPE"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			RequireEnv("DATABASE_USER"),
			RequireEnv("DATABASE_PASS"),
			RequireEnv("DATABASE_HOST"),
			RequireEnv("DATABASE_PORT"),
			RequireEnv("DATABASE_DB"),
		),
	)
	if err != nil {
		return nil, err
	}

	userRepository := NewDbUserRepository(db)

	service := NewService(userRepository)
	service.getParams = func(r *http.Request) map[string]string {
		return mux.Vars(r)
	}

	return &App{
		db:      db,
		service: service,
	}, nil
}

// Handles application shutdown (closes DB connection, etc)
func (app *App) Shutdown() {
	app.db.Close()
}

// Starts listening
func (app *App) Listen() error {
	handler := app.CreateHandler()

	return http.ListenAndServe(":80", handler)
}

// Creates and configures the router
func (app *App) CreateHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", app.service.CreateUser).Methods("POST")
	router.HandleFunc("/{userKey}", app.service.GetUser).Methods("GET")

	return router
}

// Tries to find an env var and panics if it's not found
func RequireEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s is mandatory", key))
	}

	return value
}
