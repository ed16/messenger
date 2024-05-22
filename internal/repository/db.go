package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Database setup and connection
func GetPostgresConn() (*sql.DB, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	return dbConn, err
}

// MockDB is a mock implementation of *sql.DB
type MockDB struct{}

// QueryContext is a mock implementation of QueryContext
func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	// Implement your mock behavior here
	return nil, errors.New("mock: QueryContext not implemented")
}

// ExecContext is a mock implementation of ExecContext
func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// Implement your mock behavior here
	return nil, errors.New("mock: ExecContext not implemented")
}

// PingContext is a mock implementation of PingContext
func (m *MockDB) PingContext(ctx context.Context) error {
	// Implement your mock behavior here
	return errors.New("mock: PingContext not implemented")
}

// PrepareContext is a mock implementation of PrepareContext
func (m *MockDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// Implement your mock behavior here
	return nil, errors.New("mock: PrepareContext not implemented")
}

// BeginTx is a mock implementation of BeginTx
func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// Implement your mock behavior here
	return nil, errors.New("mock: BeginTx not implemented")
}

// Close is a mock implementation of Close
func (m *MockDB) Close() error {
	// Implement your mock behavior here
	return errors.New("mock: Close not implemented")
}

// DBWithMock is a struct embedding *sql.DB and including MockDB
type DBWithMock struct {
	SQL *sql.DB // Renamed the field to avoid conflict
	MockDB
}

// GetMockDB returns a new instance of DBWithMock
func GetMockDB() *DBWithMock {
	return &DBWithMock{}
}

// DB returns the embedded *sql.DB instance
func (db *DBWithMock) DB() *sql.DB {
	return db.SQL // Returning the renamed field
}
