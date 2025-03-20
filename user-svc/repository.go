package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

// UserRepositoryInterface defines the methods that a user repository must implement
type UserRepositoryInterface interface {
	GetAllUsers(pageNum, pageSize int) ([]User, error)
	GetUserByID(id int) (User, error)
	CreateUser(name string) (User, error)
}

// UserRepository handles data access operations for users
type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
	dbType string
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB, logger *slog.Logger, dbType string) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
		dbType: dbType,
	}
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers(pageNum, pageSize int) ([]User, error) {
	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate offset
	offset := (pageNum - 1) * pageSize

	var rows *sql.Rows
	var err error

	// Use appropriate SQL syntax based on database type
	if r.dbType == "sqlite" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, name, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`, pageSize, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, name, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, pageSize, offset)
	}

	if err != nil {
		r.logger.Error("Database query failed", "error", err, "pageNum", pageNum, "pageSize", pageSize)
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
			r.logger.Error("Row scan failed", "error", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		r.logger.Error("Row iteration error", "error", err)
		return nil, err
	}

	r.logger.Info("Retrieved users", "count", len(users))
	return users, nil
}

// GetUserByID retrieves a specific user by ID
func (r *UserRepository) GetUserByID(id int) (User, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	var err error

	// Use appropriate SQL syntax based on database type
	if r.dbType == "sqlite" {
		err = r.db.QueryRowContext(ctx, `
			SELECT id, name, created_at, updated_at
			FROM users
			WHERE id = ?
		`, id).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	} else {
		// Default to PostgreSQL syntax
		err = r.db.QueryRowContext(ctx, `
			SELECT id, name, created_at, updated_at
			FROM users
			WHERE id = $1
		`, id).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Info("User not found", "id", id)
			return User{}, ErrUserNotFound
		}
		r.logger.Error("Database query failed", "error", err, "id", id)
		return User{}, err
	}

	r.logger.Info("Retrieved user", "id", id)
	return user, nil
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(name string) (User, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get current timestamp in microseconds
	now := time.Now().UnixMicro()

	var user User
	var err error

	// Use appropriate SQL syntax based on database type
	if r.dbType == "sqlite" {
		// SQLite doesn't support RETURNING, so we need to use two queries
		result, err := r.db.ExecContext(ctx, `
			INSERT INTO users (name, created_at, updated_at)
			VALUES (?, ?, ?)
		`, name, now, now)

		if err != nil {
			r.logger.Error("Failed to create user", "error", err, "name", name)
			return User{}, err
		}

		// Get the last inserted ID
		lastID, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID", "error", err)
			return User{}, err
		}

		// Fetch the created user
		err = r.db.QueryRowContext(ctx, `
			SELECT id, name, created_at, updated_at
			FROM users
			WHERE id = ?
		`, lastID).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			r.logger.Error("Failed to fetch created user", "error", err)
			return User{}, err
		}

	} else {
		// PostgreSQL supports RETURNING
		err = r.db.QueryRowContext(ctx, `
			INSERT INTO users (name, created_at, updated_at)
			VALUES ($1, $2, $3)
			RETURNING id, name, created_at, updated_at
		`, name, now, now).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	}

	if err != nil {
		r.logger.Error("Failed to create user", "error", err, "name", name)
		return User{}, err
	}

	r.logger.Info("Created user", "id", user.ID, "name", user.Name)
	return user, nil
}
