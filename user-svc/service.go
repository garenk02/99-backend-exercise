package main

import "errors"

// Define custom error types
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrRequiredField   = errors.New("required field missing")
)

// UserService handles business logic for user operations
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(pageNum, pageSize int) ([]User, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.GetAllUsers(pageNum, pageSize)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidArgument
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name string) (User, error) {
	if name == "" {
		return User{}, ErrRequiredField
	}

	return s.repo.CreateUser(name)
}
