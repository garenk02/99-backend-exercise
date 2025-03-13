package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"public-api/config"
	"public-api/domain"
	"public-api/handlers"
	"public-api/logger"
	"public-api/repository"
	"public-api/usecase"
)

func main() {
	// Setup logger
	logger.Setup()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize repositories
	userRepo := repository.NewUserRepository(cfg.UserServiceURL)
	listingRepo := repository.NewListingRepository(cfg.ListingServiceURL)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	listingUseCase := usecase.NewListingUseCase(listingRepo, userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	listingHandler := handlers.NewListingHandler(listingUseCase)

	// Setup router using standard http.ServeMux
	mux := http.NewServeMux()
	
	// Register routes
	mux.HandleFunc("/public-api/users", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		if r.Method == http.MethodPost {
			slog.Info("Request received",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			
			userHandler.CreateUser(w, r)
		} else {
			slog.Info("Method not allowed",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			
			domain.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		}
		
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
	})

	mux.HandleFunc("/public-api/listings", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		slog.Info("Request received",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		
		switch r.Method {
		case http.MethodGet:
			listingHandler.GetListings(w, r)
		case http.MethodPost:
			listingHandler.CreateListing(w, r)
		default:
			domain.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		}
		
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
	})

	// Create middleware for logging unhandled errors
	handler := logMiddleware(mux)

	// Start server
	port := cfg.ServerPort
	slog.Info("Server starting", "port", port)
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// logMiddleware logs all requests and recovers from panics
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Recovered from panic", 
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
				)
				domain.RespondWithError(w, http.StatusInternalServerError, "Internal server error", nil)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}