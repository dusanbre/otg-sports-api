package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var (
	instance *DB
	once     sync.Once
)

// DB wraps the database connection and Squirrel builder
type DB struct {
	Conn    *sql.DB
	Builder sq.StatementBuilderType
}

// GetInstance returns the singleton database instance
func GetInstance() (*DB, error) {
	var err error

	once.Do(func() {
		instance, err = newConnection()
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// newConnection creates a new database connection
func newConnection() (*DB, error) {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "otg_sports")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	// Create Squirrel builder with PostgreSQL placeholders
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &DB{
		Conn:    conn,
		Builder: builder,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
