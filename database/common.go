package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// SQL Related Constants (moved from logindb.go)
const (
	// User Management Queries
	SQL_INSERT_USER             = `INSERT INTO users (username, password_hash) VALUES (?, ?)` // Changed to password_hash
	SQL_CHECK_USER_EXISTS       = `SELECT COUNT(*) FROM users WHERE username = ?`
	SQL_SELECT_USER_PASSWORD    = `SELECT password_hash FROM users WHERE username = ?`               // Changed to password_hash
	SQL_SELECT_USER_BY_USERNAME = `SELECT id, username, password_hash FROM users WHERE username = ?` // Added password_hash for login flow
	SQL_SELECT_USER_BY_ID       = `SELECT id, username FROM users WHERE id = ?`
	SQL_UPDATE_USER_PASSWORD    = `UPDATE users SET password_hash = ? WHERE username = ?` // Changed to password_hash
	// Removed SQL_LOGIN - direct plaintext password query is bad
)

// Session Management Queries (from sessiondb.go)
const (
	SQL_INSERT_SESSION            = `INSERT INTO sessions (sessions_id, user_id, expires_at) VALUES (?, ?, ?)`   // Changed to user_id
	SQL_SELECT_SESSION_BY_ID      = `SELECT user_id, created_at, expires_at FROM sessions WHERE sessions_id = ?` // Changed to user_id
	SQL_DELETE_SESSION_BY_ID      = `DELETE FROM sessions WHERE sessions_id = ?`
	SQL_DELETE_EXPIRED_SESSIONS   = `DELETE FROM sessions WHERE expires_at < NOW()`
	SQL_UPDATE_SESSION_EXPIRATION = `UPDATE sessions SET expires_at = ? WHERE sessions_id = ?`
	// If you want to list/delete sessions by username, you'd need a join or another lookup
	// SQL_SELECT_SESSIONS_BY_USERNAME = `SELECT s.session_id, s.created_at, s.expires_at FROM sessions s JOIN users u ON s.user_id = u.id WHERE u.username = ?`
	// SQL_DELETE_ALL_SESSIONS_FOR_USER = `DELETE FROM sessions WHERE user_id IN (SELECT id FROM users WHERE username = ?)`
)

// primary type for interacting with Database (renamed to avoid conflict if you have multiple DBs)
type DBService struct {
	db *sql.DB
}

// helper type for dealing with user databases
type User struct { // Renamed from Person to User
	ID           int // Changed to int to match your SQL_SELECT_USER_BY_ID. If UUIDs, use string.
	Username     string
	PasswordHash string // Store the hashed password here
}

const (
	DriverName            string = "mysql"
	UserSessionDataSource string = "ethan:040323@tcp(10.0.0.173:3306)user_session_db/?parseTime=true"
	CONNECTIONS           int    = 50
)

// NewDBService creates a DBService object pointer with a database connection.
// It uses a context for the initial ping.
func NewDBService(ctx context.Context, dataSourceName string) (*DBService, error) {
	db, err := sql.Open(DriverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection for %s: %w", dataSourceName, err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(CONNECTIONS)
	db.SetMaxIdleConns(CONNECTIONS)

	err = db.PingContext(ctx)
	if err != nil {
		db.Close() // Close connection if ping fails
		return nil, fmt.Errorf("error pinging database %s: %w", dataSourceName, err)
	}

	fmt.Printf("Connection established to database: %s\n", dataSourceName)

	return &DBService{db: db}, nil
}

// closes database connection associated with a Service Object
func (s *DBService) Close() error {
	log.Println("Closing database connection.")
	return s.db.Close()
}
