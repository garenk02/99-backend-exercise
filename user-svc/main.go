package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Set up structured logging with slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("Error loading .env file, using environment variables", "error", err)
	}

	// Get database type (postgres or sqlite)
	dbType := getEnv("DB_TYPE", "postgres")
	dbType = strings.ToLower(dbType)

	var db *sql.DB
	var err error

	// Connect to appropriate database based on DB_TYPE
	if dbType == "sqlite" {
		// SQLite connection
		sqlitePath := getEnv("SQLITE_DB_PATH", "./userservice.db")
		slog.Info("Using SQLite database", "path", sqlitePath)

		db, err = sql.Open("sqlite3", sqlitePath)
		if err != nil {
			slog.Error("Failed to connect to SQLite database", "error", err)
			os.Exit(1)
		}
	} else {
		// Default to PostgreSQL
		dbConfig := getDBConfig()
		slog.Info("Using PostgreSQL database", "host", dbConfig.Host, "dbname", dbConfig.Name)

		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			// Try SQLite as fallback
			slog.Warn("Failed to connect to PostgreSQL, trying SQLite as fallback", "error", err)

			sqlitePath := getEnv("SQLITE_DB_PATH", "./userservice.db")
			db, err = sql.Open("sqlite3", sqlitePath)
			if err != nil {
				slog.Error("Failed to connect to SQLite fallback", "error", err)
				os.Exit(1)
			}
			slog.Info("Connected to SQLite fallback database", "path", sqlitePath)
			dbType = "sqlite" // Update type for table creation
		}
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err, "dbType", dbType)
		os.Exit(1)
	}
	slog.Info("Successfully connected to database", "type", dbType)

	// Create user table if it doesn't exist (with DB-specific SQL)
	if err := initDB(db, dbType); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	serverPort := getServerPort()

	// Initialize repository and service with database type
	userRepo := NewUserRepository(db, logger, dbType)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	// Start server
	slog.Info("Server starting", "port", serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), mux); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// getDBConfig reads database configuration from environment variables
func getDBConfig() DBConfig {
	// Get database configuration with defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPortStr := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "userservice")

	// Parse port number
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		slog.Warn("Invalid DB_PORT, using default", "value", dbPortStr, "default", 5432)
		dbPort = 5432
	}

	return DBConfig{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
	}
}

// getServerPort reads server port from environment variables
func getServerPort() int {
	portStr := getEnv("SERVER_PORT", "6001")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Warn("Invalid SERVER_PORT, using default", "value", portStr, "default", 6001)
		return 6001
	}
	return port
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// initDB creates the database schema based on the database type
func initDB(db *sql.DB, dbType string) error {
	var createTableSQL string

	// Use appropriate SQL syntax based on database type
	if dbType == "sqlite" {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				created_at INTEGER NOT NULL,
				updated_at INTEGER NOT NULL
			)
		`
	} else {
		// Default to PostgreSQL syntax
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL
			)
		`
	}

	_, err := db.Exec(createTableSQL)
	return err
}
