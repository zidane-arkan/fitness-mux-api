// main_test.go
package main

import (
	"log"
	"os"
	"testing"
)

var a App

// TestMain is the main test function for the package. It initializes the application,
// checks if the required table exists, runs the tests, clears the table, and exits.
func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)
	checkTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

// checkTableExists checks if a table exists in the database.
// If the table does not exist, it creates the table using the createTableQuery.
func checkTableExists() {
	if _, err := a.DB.Exec(createTableQuery); err != nil {
		log.Fatal(err)
	}
}

// clearTable deletes all records from the exercises table and resets the ID sequence.
func clearTable() {
	a.DB.Exec("DELETE FROM exercises")
	a.DB.Exec("ALTER SEQUENCE exercises_id_seq RESTART WITH 1")
}

// createTableQuery is a SQL query used to create the "exercises" table if it doesn't already exist.
const createTableQuery = `CREATE TABLE IF NOT EXISTS exercises
(
	id SERIAL,
	name TEXT NOT NULL,
	workoutType TEXT NOT NULL,
	sets INTEGER,
	CONSTRAINT exercise_pkey PRIMARY KEY (id)	
)`
