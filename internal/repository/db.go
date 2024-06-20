package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// GetPostgresConn attempts to connect to the PostgreSQL database with retries.
func GetPostgresConn() (*sql.DB, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var dbConn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		dbConn, err = sql.Open("postgres", dsn)
		if err != nil {
			fmt.Printf("Failed to open connection to database: %v. Retrying in 2 seconds... (%d/10)\n", err, i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		err = dbConn.Ping()
		if err != nil {
			fmt.Printf("Failed to ping database: %v. Retrying in 2 seconds... (%d/10)\n", err, i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		// Connection successful
		return dbConn, nil
	}

	// All attempts failed
	return nil, fmt.Errorf("could not connect to the database after several attempts: %v", err)
}
