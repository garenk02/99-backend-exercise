// handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"public-api/domain"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		domain.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if request.Name == "" {
		domain.RespondWithError(w, http.StatusBadRequest, "Name is required", nil)
		return
	}

	// Create user
	user, err := h.userUseCase.CreateUser(request.Name)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Log success
	slog.Info("User created successfully", "user_id", user.ID, "name", user.Name)

	// Prepare response
	response := struct {
		User *domain.User `json:"user"`
	}{
		User: user,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}