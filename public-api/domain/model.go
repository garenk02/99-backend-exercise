package domain

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Listing struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int    `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	User        *User  `json:"user,omitempty"`
}

// ListingWithUser represents a listing with embedded user data
type ListingWithUser struct {
	Listing
	User User `json:"user"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetUserByID(id int) (*User, error)
	GetUsers(pageNum, pageSize int) ([]*User, error)
	CreateUser(name string) (*User, error)
}

// ListingRepository defines the interface for listing data operations
type ListingRepository interface {
	GetListings(pageNum, pageSize int, userID *int) ([]*Listing, error)
	CreateListing(userID int, listingType string, price int) (*Listing, error)
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	GetUserByID(id int) (*User, error)
	GetUsers(pageNum, pageSize int) ([]*User, error)
	CreateUser(name string) (*User, error)
}

// ListingUseCase defines the interface for listing business logic
type ListingUseCase interface {
	GetListings(pageNum, pageSize int, userID *int) ([]*ListingWithUser, error)
	CreateListing(userID int, listingType string, price int) (*Listing, error)
}
