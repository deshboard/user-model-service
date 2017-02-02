package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Integration test initialization
func integrationSetUp() {
	setupDatabase()
}

// Cleanup after integration tests
func integrationTearDown() {
	teardownDatabase()
}

func setupDatabase() {
	db, _ = sqlx.Open(
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

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	cleanupDatabase()
}

func teardownDatabase() {
	cleanupDatabase()

	db.Close()
}

func cleanupDatabase() {
	db.MustExec("DELETE FROM users")
}
