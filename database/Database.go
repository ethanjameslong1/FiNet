package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// SQL Related Constants
const (
	//db.ExecContext
	SQL_INSERT_USER = `INSERT INTO users (username, password) VALUES (?, ?)`

	//db.QueryRowContext
	SQL_CHECK_USER_EXISTS       = `SELECT COUNT(*) FROM users WHERE username = ?`
	SQL_SELECT_USER_PASSWORD    = `SELECT password FROM users WHERE username = ?`
	SQL_SELECT_USER_BY_USERNAME = `SELECT id, username, created_at FROM users WHERE username = ?`
	SQL_SELECT_USER_BY_ID       = `SELECT id, username, created_at FROM users WHERE id = ?`
	SQL_UPDATE_USER_PASSWORD    = `UPDATE users SET password = ? WHERE username = ?`
)

type Service struct {
	db *sql.DB
}

type Person struct {
	Username  string
	Password  string
	Id        int
	CreatedAt time.Time
}

var DriverName string = "mysql" //so I don't forget the Driver i'm using teehee

// creates a Service object pointer with a database connection, requires a driver and datasource location
func NewService(driverName, dataSourceName string) (*Service, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Println("Connection established to database")

	return &Service{db: db}, nil
}

// closes database connection associated with a Service Object
func (s *Service) Close() error {
	log.Println("Closing database connection.")
	return s.db.Close()
}

// adds a user to the database associated with a Service Object
func AddUser(ctx context.Context, s *Service, name string, password string) (bool, error) {
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

func QueryUserByName(ctx context.Context, s *Service, name string) (Person, error) {
	person := Person{}
	row := s.db.QueryRowContext(ctx, SQL_SELECT_USER_BY_USERNAME, name)
	err := row.Scan(&person.Id, &person.Username, &person.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Person{}, fmt.Errorf("user not found: %w", err)
		}
		return Person{}, fmt.Errorf("Error finding User: %w", err)
	}
	return person, nil
}
