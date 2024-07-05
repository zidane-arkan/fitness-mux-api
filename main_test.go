// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

// TestEmptyTable tests the behavior when the table is empty.
func TestEmptyTable(t *testing.T) {
	clearTable()
	// Membuat permintaan HTTP GET ke endpoint "/exercises".
	req, _ := http.NewRequest("GET", "/exercises", nil)
	// Menjalankan permintaan HTTP yang dibuat dan mengembalikan respons
	response := executeReq(req)
	// Memeriksa apakah kode respons yang diterima adalah 200 OK.
	checkResponseCode(t, http.StatusOK, response.Code)
	// Memeriksa apakah tubuh respons adalah array kosong ("[]"). Jika tidak, mencetak pesan kesalahan.
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected empty array. Got %s", body)
	}
}

// TestGetNonExistentExercise tests the behavior when trying to get a non-existent exercise.
func TestGetNonExistentExercise(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/exercise/11", nil)
	response := executeReq(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Exercise not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Exercise not found'. Got '%s'", m["error"])
	}
}

// TestCreateExercise is a unit test function that tests the creation of an exercise.
// It sends a POST request to the "/exercise" endpoint with a JSON payload representing a workout.
// The function checks if the exercise is created successfully by verifying the response body.
// It asserts that the exercise name, workout type, sets, and ID match the expected values.
func TestCreateExercise(t *testing.T) {
	clearTable()

	// jsonStr is a byte slice containing a JSON string representing a workout.
	var jsonStr = []byte(`{"name": "Squat", "workoutType": "strength", "sets": 4}`)
	// Create a new HTTP request with the specified method, URL path, and request body.
	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeReq(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	// m is a map that stores key-value pairs where the keys are of type string and the values are of type interface{}.
	var m map[string]interface{}
	// Fungsi json.Unmarshal mengonversi data JSON dari response.Body.Bytes() menjadi peta (map) m.
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Squat" {
		t.Errorf("Expected exercise name is 'Squat' but get %s", m["name"])
	}

	if m["workoutType"] != "strength" {
		t.Errorf("Expected exercise workoutType is 'Strength' but get %s", m["workoutType"])
	}

	if m["sets"] != 4 {
		t.Errorf("Expected exercise sets is '4' but get %d", m["sets"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected exercise id is '1 but get %v", m["id"])
	}
}

// TestGetExercise is a unit test function that tests the GetExercise handler.
func TestGetExercise(t *testing.T) {
	clearTable()
	addExercises(1)

	req, _ := http.NewRequest("GET", "/exercise/1", nil)
	response := executeReq(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
func addExercises(count int) {
	if count > 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec(
			"INSERT INTO exercises(name, workoutType, sets) VALUES(%s, %s, %d)",
			"Name "+strconv.Itoa(i), "workout "+strconv.Itoa(i), (i+1)*2)

	}
}

// This test begins by adding a exercise to the database directly.
// It then uses the end point to update this record with new details.
// We finally test the following things:
// - That the status code is 200, indicating success, and
// - That the response contains the JSON representation of the exercise with the updated details.
func TestUpdateExercise(t *testing.T) {
	clearTable()
	addExercises(1)

	req, _ := http.NewRequest("GET", "/exercise/1", nil)
	response := executeReq(req)

	var originalExercise map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalExercise)

	var jsonStr = []byte(`{{"name": "Updated BP", "workoutType": "Update strength", "sets": 5}}`)
	req, _ = http.NewRequest("PUT", "/exercise/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeReq(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalExercise["id"] {
		t.Errorf("Expected exercise id to remain the same %v. Got %v", originalExercise["id"], m["id"])
	}

	if m["name"] == originalExercise["name"] {
		t.Errorf("Expected exercise name to change from %s to %s but remain the same. Got %s", originalExercise["name"], m["name"], m["name"])
	}

	if m["workoutType"] == originalExercise["workoutType"] {
		t.Errorf("Expected exercise workoutType to change from %s to %s. Got %s", originalExercise["workoutType"], m["workoutType"], m["workoutType"])
	}

	if m["sets"] == originalExercise["sets"] {
		t.Errorf("Expected exercise sets to change from %d to %d. Got %d", originalExercise["sets"], m["sets"], m["sets"])
	}
}

// In this test, we first create a product and test that it exists. We then use the endpoint to delete the product. Finally we try to access the product at the appropriate endpoint and test that it doesn’t exist.
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addExercises(1)

	req, _ := http.NewRequest("GET", "/exercise/1", nil)
	response := executeReq(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/exercise/1", nil)
	response = executeReq(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/exercise/1", nil)
	response = executeReq(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// executeReq executes the given HTTP request and returns the response recorder.
func executeReq(req *http.Request) *httptest.ResponseRecorder {
	//  Membuat perekam respons untuk menangkap respons dari server.
	rr := httptest.NewRecorder()
	// Menjalankan permintaan HTTP melalui router aplikasi dan menangkap respons.
	a.Router.ServeHTTP(rr, req)
	// Mengembalikan respons yang telah direkam.
	return rr
}

// checkResponseCode checks if the expected response code matches the actual response code.
// If they don't match, it reports an error using the testing.T.Errorf function.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d \n", expected, actual)
	}
}
