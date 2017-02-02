package main

import (
	"fmt"
	"net/http"

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
// Make sure the process does not exit before this is called
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

	router.HandleFunc("/_status/healthz", app.HealthStatus).Methods("GET")
	router.HandleFunc("/_status/readiness", app.ReadinessStatus).Methods("GET")

	return router
}

// Checks if the app is up and running
func (app *App) HealthStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Checks if the app is ready for accepting request (eg. database is available as well)
func (app *App) ReadinessStatus(w http.ResponseWriter, r *http.Request) {
	if err := app.db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("error"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
