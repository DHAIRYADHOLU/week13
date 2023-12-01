package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentTimeInToronto(t *testing.T) {
	// Call the function to get current time in Toronto
	torontoTime, err := getCurrentTimeInToronto()

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the returned time is not zero
	assert.NotEqual(t, time.Time{}, torontoTime, "Expected non-zero time, got zero time")
}

func TestLogTimeToDatabase(t *testing.T) {
	// Initialize a temporary in-memory database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Call the function to log time to the database
	currentTime := time.Now()
	err = logTimeToDatabase(db, currentTime)

	// Assert that there is no error
	assert.NoError(t, err)

	// Query the database to check if the time was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM time_log WHERE timestamp = ?", currentTime).Scan(&count)

	// Assert that the count is 1, indicating that the time was inserted
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected count to be 1, got %d", count)
}

func TestGetCurrentTimeEndpoint(t *testing.T) {
	// Initialize the Gin router
	router := gin.Default()

	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make a request to the /current-time endpoint
		r.URL.Path = "/current-time"
		router.ServeHTTP(w, r)
	}))
	defer ts.Close()

	// Make a GET request to the test server
	res, err := http.Get(ts.URL + "/current-time")
	assert.NoError(t, err)
	defer res.Body.Close()

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
