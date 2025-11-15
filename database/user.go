package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func (s *DBService) AddUser(ctx context.Context, name string, password string) (bool, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}
	r, err := s.db.ExecContext(ctx, SQL_INSERT_USER, name, string(hashedPW))
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

func (s *DBService) GetUserByName(ctx context.Context, username string) (User, error) {
	person := User{}
	row := s.db.QueryRowContext(ctx, SQL_SELECT_USER_BY_USERNAME, username)
	err := row.Scan(&person.ID, &person.Username, &person.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("User not found: %w", err)
		}
		return User{}, fmt.Errorf("Error finding User: %w", err)
	}
	return person, nil
}

func (s *DBService) GetUserByID(ctx context.Context, ID int) (User, error) {
	person := User{}
	row := s.db.QueryRowContext(ctx, SQL_SELECT_USER_BY_ID, ID)
	err := row.Scan(&person.ID, &person.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("User not found: %w", err)
		}
		return User{}, fmt.Errorf("Error finding User: %w", err)
	}
	return person, nil
}

func (s *DBService) AuthenticateUser(ctx context.Context, name string, passwordHash string) (User, error) {
	user, err := s.GetUserByName(ctx, name)
	if err != nil {
		return User{}, fmt.Errorf("Authentication failed: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return User{}, fmt.Errorf("Authentication failed: %w", err)
		}
		return User{}, fmt.Errorf("Authentication failed: bcrypt error: %w", err)
	}
	user.PasswordHash = ""
	return user, nil

}
