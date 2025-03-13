// repository/user_repository.go
package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"public-api/domain"
	"strings"
)

type UserRepository struct {
	baseURL string
}

func NewUserRepository(baseURL string) *UserRepository {
	return &UserRepository{
		baseURL: baseURL,
	}
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	slog.Debug("Fetching user by ID", "user_id", id)

	// Fetch users and filter by ID
	users, err := r.GetUsers(1, 100)
	if err != nil {
		slog.Error("Failed to fetch users", "error", err)
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	slog.Warn("User not found", "user_id", id)
	return nil, errors.New("user not found")
}

func (r *UserRepository) GetUsers(pageNum, pageSize int) ([]*domain.User, error) {
	slog.Debug("Fetching users", "page_num", pageNum, "page_size", pageSize)

	// Build URL with query parameters
	reqURL := fmt.Sprintf("%s/users?page_num=%d&page_size=%d", r.baseURL, pageNum, pageSize)

	// Make HTTP request
	slog.Debug("Making request to user service", "url", reqURL)
	resp, err := http.Get(reqURL)
	if err != nil {
		slog.Error("Error making request to user service", "error", err)
		return nil, fmt.Errorf("error making request to user service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		slog.Error("User service returned non-200 status", "status", resp.StatusCode)
		return nil, fmt.Errorf("user service returned non-200 status: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Result bool `json:"result"`
		Users  []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding response from user service", "error", err)
		return nil, fmt.Errorf("error decoding response from user service: %w", err)
	}

	// Convert to domain model
	users := make([]*domain.User, 0, len(response.Users))
	for _, u := range response.Users {
		users = append(users, &domain.User{
			ID:        u.ID,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	slog.Debug("Fetched users successfully", "count", len(users))
	return users, nil
}

func (r *UserRepository) CreateUser(name string) (*domain.User, error) {
	slog.Debug("Creating user", "name", name)

	// Prepare form data
	data := url.Values{}
	data.Set("name", name)

	// Make HTTP request
	reqURL := fmt.Sprintf("%s/users", r.baseURL)
	slog.Debug("Making request to user service", "url", reqURL)

	resp, err := http.Post(
		reqURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		slog.Error("Error making request to user service", "error", err)
		return nil, fmt.Errorf("error making request to user service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		slog.Error("User service returned status", "status", resp.StatusCode)
		return nil, fmt.Errorf("user service returned status: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Result bool `json:"result"`
		User   struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding response from user service", "error", err)
		return nil, fmt.Errorf("error decoding response from user service: %w", err)
	}

	// Convert to domain model
	user := &domain.User{
		ID:        response.User.ID,
		Name:      response.User.Name,
		CreatedAt: response.User.CreatedAt,
		UpdatedAt: response.User.UpdatedAt,
	}

	slog.Debug("User created successfully", "user_id", user.ID)
	return user, nil
}
