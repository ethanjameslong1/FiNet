package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (s *DBService) AddSession(ctx context.Context, sessionID uuid.UUID, userID int, duration time.Duration) (bool, error) {
	expiresAt := time.Now().Add(duration)
	r, err := s.db.ExecContext(ctx, SQL_INSERT_SESSION, sessionID.String(), userID, expiresAt)
	if err != nil {
		return false, fmt.Errorf("error inserting session: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error finding rows impacted: %w", err)
	}
	if rows != 1 {
		return false, fmt.Errorf("expected 1 row impacted for session insertion, got %d", rows)
	}
	return true, nil
}

func (s *DBService) GetSessionByID(ctx context.Context, sessionID string) (Session, error) {
	session := Session{}
	row := s.db.QueryRowContext(ctx, SQL_SELECT_SESSION_BY_ID, sessionID)
	err := row.Scan(&session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Session{}, fmt.Errorf("session not found: %w", err)
		}
		return Session{}, fmt.Errorf("error retrieving session by ID: %w", err)
	}
	session.ID = sessionID
	return session, nil
}

func (s *DBService) DeleteSessionByID(ctx context.Context, sessionID string) (bool, error) {
	r, err := s.db.ExecContext(ctx, SQL_DELETE_SESSION_BY_ID, sessionID)
	if err != nil {
		return false, fmt.Errorf("error deleting session: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking rows affected for session delete: %w", err)
	}
	return rows > 0, nil
}

func (s *DBService) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	r, err := s.db.ExecContext(ctx, SQL_DELETE_EXPIRED_SESSIONS)
	if err != nil {
		return 0, fmt.Errorf("error deleting expired sessions: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking rows affected for expired session cleanup: %w", err)
	}
	return rows, nil
}

func (s *DBService) UpdateSessionExpiration(ctx context.Context, sessionID string, newDuration time.Duration) (bool, error) {
	newExpiresAt := time.Now().Add(newDuration)
	r, err := s.db.ExecContext(ctx, SQL_UPDATE_SESSION_EXPIRATION, newExpiresAt, sessionID)
	if err != nil {
		return false, fmt.Errorf("error updating session expiration: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking rows affected for session expiration update: %w", err)
	}
	return rows > 0, nil
}
