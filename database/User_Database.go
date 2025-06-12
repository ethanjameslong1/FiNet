package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const COMPILE_VERSION = "20250609_V_FINAL_FIX"

// SQL Related Constants
const (
	//db.ExecContext
	SQL_INSERT_USER = `INSERT INTO users (username, password) VALUES (?, ?)`

	//db.QueryRowContext
	SQL_CHECK_USER_EXISTS       = `SELECT COUNT(*) FROM users WHERE username = ?`
	SQL_SELECT_USER_PASSWORD    = `SELECT password FROM users WHERE username = ?`
	SQL_SELECT_USER_BY_USERNAME = `SELECT id, username FROM users WHERE username = ?`
	SQL_SELECT_USER_BY_ID       = `SELECT id, username FROM users WHERE id = ?`
	SQL_UPDATE_USER_PASSWORD    = `UPDATE users SET password = ? WHERE username = ?`

	SQL_LOGIN = `SELECT id, username FROM users WHERE username = ? AND password = ?`
)

// primary type for interacting with Database
type Service struct {
	db *sql.DB
}

// helper type for dealing with user databases
type Person struct {
	Username string
	Password string
	Id       int
}

const (
	DriverName  string = "mysql"                                     //so I don't forget the Driver i'm using teehee
	DataSource  string = "ethan:040323@tcp(10.0.0.89:3306)/my_go_db" //see prior comment... teehee
	CONNECTIONS int    = 50
)

// creates a Service object pointer with a database connection, requires a driver and datasource location, DriverName and DataSource constants
func NewService(driverName, dataSourceName string) (*Service, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(CONNECTIONS)
	db.SetMaxIdleConns(CONNECTIONS)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Println("Connection established to database")
	fmt.Printf("DEBUG: Database package compiled with version: %s\n", COMPILE_VERSION)

	return &Service{db: db}, nil
}

// closes database connection associated with a Service Object
func (s *Service) Close() error {
	log.Println("Closing database connection.")
	return s.db.Close()
}

// adds a user to the database given creation context, name, and string
// currently not in use, don't plan on adding functionality to add users easily. For my current intents only the one user I added manually is neccesary
func (s *Service) AddUser(ctx context.Context, name string, password string) (bool, error) {
	r, err := s.db.ExecContext(ctx, SQL_INSERT_USER, name, password)
	if err != nil {
		return false, fmt.Errorf("Error inserting User: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Error finding Rows Affected: %w", err)
	}
	if rows != 1 {
		return false, fmt.Errorf("Expected rows impacted to be 1, rows impacted %d", rows)
	}
	return true, nil
}

// returns a Person type that matches query name
func (s *Service) QueryUserByName(ctx context.Context, name string) (Person, error) {
	person := Person{}
	row := s.db.QueryRowContext(ctx, SQL_SELECT_USER_BY_USERNAME, name)
	err := row.Scan(&person.Id, &person.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return Person{}, fmt.Errorf("User not found: %w", err)
		}
		return Person{}, fmt.Errorf("Error finding User: %w", err)
	}
	return person, nil
}

// checks name and password and returns matching Person object
func (s *Service) LoginQuery(ctx context.Context, name string, password string) (Person, error) {
	person := Person{}
	var UCount int
	err := s.db.QueryRowContext(ctx, SQL_CHECK_USER_EXISTS, name).Scan(&UCount)
	if err != nil {
		return Person{}, fmt.Errorf("Error Checking User Existence: %w", err)
	}
	if UCount == 0 {
		return Person{}, fmt.Errorf("Invalid Username or Password: user %s does not exist", name)
	}

	row := s.db.QueryRowContext(ctx, SQL_LOGIN, name, password)

	err = row.Scan(&person.Id, &person.Username)
	if err != nil {
		return Person{}, fmt.Errorf("Error Logging in: %w", err)
	}
	if person.Username != "" {
		return person, nil //successful login
	} else {
		return Person{}, fmt.Errorf("Invalid Username or Password: %w", err) //wrong password
	}

}
