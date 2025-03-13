package main

// User represents a user in the system
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// UsersResponse is the response format for GET /users
type UsersResponse struct {
	Result bool   `json:"result"`
	Users  []User `json:"users"`
}

// UserResponse is the response format for GET /users/{id} and POST /users
type UserResponse struct {
	Result bool `json:"result"`
	User   User `json:"user"`
}
