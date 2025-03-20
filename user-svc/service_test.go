package main

import (
	"errors"
	"reflect"
	"testing"
)

// MockUserRepository implements UserRepository methods for testing
type MockUserRepository struct {
	getUserByIDFn func(id int) (User, error)
	getAllUsersFn func(pageNum, pageSize int) ([]User, error)
	createUserFn  func(name string) (User, error)
}

// GetUserByID mocks the repository method
func (m *MockUserRepository) GetUserByID(id int) (User, error) {
	return m.getUserByIDFn(id)
}

// GetAllUsers mocks the repository method
func (m *MockUserRepository) GetAllUsers(pageNum, pageSize int) ([]User, error) {
	return m.getAllUsersFn(pageNum, pageSize)
}

// CreateUser mocks the repository method
func (m *MockUserRepository) CreateUser(name string) (User, error) {
	return m.createUserFn(name)
}

// Setup test data
func setupTestUsers() []User {
	return []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
}

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	if service.repo != mockRepo {
		t.Errorf("Expected repo to be %v, got %v", mockRepo, service.repo)
	}
}

func TestGetAllUsers(t *testing.T) {
	testUsers := setupTestUsers()

	tests := []struct {
		name        string
		pageNum     int
		pageSize    int
		expected    []User
		expectedErr error
		mockFn      func(pageNum, pageSize int) ([]User, error)
	}{
		{
			name:        "Valid pagination",
			pageNum:     1,
			pageSize:    10,
			expected:    testUsers,
			expectedErr: nil,
			mockFn: func(pageNum, pageSize int) ([]User, error) {
				return testUsers, nil
			},
		},
		{
			name:        "Invalid page number defaults to 1",
			pageNum:     0,
			pageSize:    10,
			expected:    testUsers,
			expectedErr: nil,
			mockFn: func(pageNum, pageSize int) ([]User, error) {
				if pageNum != 1 {
					t.Errorf("Expected pageNum to default to 1, got %d", pageNum)
				}
				return testUsers, nil
			},
		},
		{
			name:        "Invalid page size defaults to 10",
			pageNum:     1,
			pageSize:    0,
			expected:    testUsers,
			expectedErr: nil,
			mockFn: func(pageNum, pageSize int) ([]User, error) {
				if pageSize != 10 {
					t.Errorf("Expected pageSize to default to 10, got %d", pageSize)
				}
				return testUsers, nil
			},
		},
		{
			name:        "Repository error",
			pageNum:     1,
			pageSize:    10,
			expected:    nil,
			expectedErr: errors.New("database error"),
			mockFn: func(pageNum, pageSize int) ([]User, error) {
				return nil, errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				getAllUsersFn: tt.mockFn,
			}
			service := NewUserService(mockRepo)

			users, err := service.GetAllUsers(tt.pageNum, tt.pageSize)

			if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(users, tt.expected) {
				t.Errorf("Expected users %v, got %v", tt.expected, users)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testUsers := setupTestUsers()

	tests := []struct {
		name        string
		id          int
		expected    User
		expectedErr error
		mockFn      func(id int) (User, error)
	}{
		{
			name:        "Valid ID",
			id:          1,
			expected:    testUsers[0],
			expectedErr: nil,
			mockFn: func(id int) (User, error) {
				return testUsers[0], nil
			},
		},
		{
			name:        "Invalid ID",
			id:          0,
			expected:    User{},
			expectedErr: ErrInvalidArgument,
			mockFn: func(id int) (User, error) {
				t.Errorf("Mock should not be called with invalid ID")
				return User{}, nil
			},
		},
		{
			name:        "User not found",
			id:          999,
			expected:    User{},
			expectedErr: ErrUserNotFound,
			mockFn: func(id int) (User, error) {
				return User{}, errors.New("user not found")
			},
		},
		{
			name:        "Repository error",
			id:          1,
			expected:    User{},
			expectedErr: errors.New("database error"),
			mockFn: func(id int) (User, error) {
				return User{}, errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				getUserByIDFn: tt.mockFn,
			}
			service := NewUserService(mockRepo)

			user, err := service.GetUserByID(tt.id)

			if err != tt.expectedErr && (err == nil || tt.expectedErr == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(user, tt.expected) {
				t.Errorf("Expected user %v, got %v", tt.expected, user)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	testUser := User{ID: 1, Name: "Alice"}

	tests := []struct {
		name        string
		userName    string
		expected    User
		expectedErr error
		mockFn      func(name string) (User, error)
	}{
		{
			name:        "Valid name",
			userName:    "Alice",
			expected:    testUser,
			expectedErr: nil,
			mockFn: func(name string) (User, error) {
				return testUser, nil
			},
		},
		{
			name:        "Empty name",
			userName:    "",
			expected:    User{},
			expectedErr: ErrRequiredField,
			mockFn: func(name string) (User, error) {
				t.Errorf("Mock should not be called with empty name")
				return User{}, nil
			},
		},
		{
			name:        "Repository error",
			userName:    "Alice",
			expected:    User{},
			expectedErr: errors.New("database error"),
			mockFn: func(name string) (User, error) {
				return User{}, errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				createUserFn: tt.mockFn,
			}
			service := NewUserService(mockRepo)

			user, err := service.CreateUser(tt.userName)

			if err != tt.expectedErr && (err == nil || tt.expectedErr == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(user, tt.expected) {
				t.Errorf("Expected user %v, got %v", tt.expected, user)
			}
		})
	}
}
