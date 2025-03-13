package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetAllUsers handles GET /users request
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse page parameters
	pageNum := 1
	pageSize := 10

	pageNumStr := r.URL.Query().Get("page_num")
	if pageNumStr != "" {
		if num, err := strconv.Atoi(pageNumStr); err == nil && num > 0 {
			pageNum = num
		}
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	if pageSizeStr != "" {
		if size, err := strconv.Atoi(pageSizeStr); err == nil && size > 0 {
			pageSize = size
		}
	}

	// Get users
	users, err := h.service.GetAllUsers(pageNum, pageSize)
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	response := UsersResponse{
		Result: true,
		Users:  users,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUser handles GET /users/{id} request
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// http.Error(w, "User not found", http.StatusNotFound)
			response := UserResponse{
				Result: false,
				User:   User{},
			}

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create response
	response := UserResponse{
		Result: true,
		User:   user,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUser handles POST /users request
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get name parameter
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.service.CreateUser(name)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	response := UserResponse{
		Result: true,
		User:   user,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
