package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Database connection parameters
const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "toronto_time_db"
)

func main() {
	// Initialize the database connection
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Initialize the Gin router
	router := gin.Default()

	// Define API endpoint to get current time in Toronto
	router.GET("/current-time", func(c *gin.Context) {
		// Get current time in Toronto's timezone
		torontoTime, err := getCurrentTimeInToronto()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Toronto time"})
			return
		}

		// Log current time to the database
		err = logTimeToDatabase(db, torontoTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log time to database"})
			return
		}

		// Respond with the current time in JSON format
		c.JSON(http.StatusOK, gin.H{"toronto_time": torontoTime.Format(time.RFC3339)})
	})

	// Run the server on port 9090
	err = router.Run(":9090")
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

// getCurrentTimeInToronto returns the current time adjusted to Toronto's timezone
func getCurrentTimeInToronto() (time.Time, error) {
	// Specify the timezone for Toronto
	torontoLocation, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return time.Time{}, err
	}

	// Get the current time in Toronto's timezone
	return time.Now().In(torontoLocation), nil
}

// logTimeToDatabase inserts the current time into the time_log table
func logTimeToDatabase(db *sql.DB, currentTime time.Time) error {
	// Prepare the SQL statement for inserting into the time_log table
	stmt, err := db.Prepare("INSERT INTO time_log (timestamp) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the current time
	_, err = stmt.Exec(currentTime)
	return err
}
