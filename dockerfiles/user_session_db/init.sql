CREATE DATABASE IF NOT EXISTS user_session_db;
CREATE USER IF NOT EXISTS 'ethan'@'%' IDENTIFIED BY '040323';
-- init.sql
-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the sessions table
CREATE TABLE IF NOT EXISTS sessions (
    sessions_id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id INT,
    expires_at DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Optional: Add some initial data for testing
INSERT INTO users (username, password_hash, email) VALUES
('testuser', 'hashed_password_123', 'test@example.com');
