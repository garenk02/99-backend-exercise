package usecase

import (
	"fmt"
	"public-api/domain"
)

type ListingUseCase struct {
	listingRepo domain.ListingRepository
	userRepo    domain.UserRepository
}

func NewListingUseCase(listingRepo domain.ListingRepository, userRepo domain.UserRepository) *ListingUseCase {
	return &ListingUseCase{
		listingRepo: listingRepo,
		userRepo:    userRepo,
	}
}

func (u *ListingUseCase) GetListings(pageNum, pageSize int, userID *int) ([]*domain.ListingWithUser, error) {
	// Get listings
	listings, err := u.listingRepo.GetListings(pageNum, pageSize, userID)
	if err != nil {
		return nil, err
	}

	// Create a map to store unique user IDs
	userIDs := make(map[int]bool)
	for _, listing := range listings {
		userIDs[listing.UserID] = true
	}

	// Get users for all listings
	users := make(map[int]*domain.User)
	for userID := range userIDs {
		user, err := u.userRepo.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("error fetching user data for listing: %w", err)
		}
		users[userID] = user
	}

	// Combine listings with user data
	listingsWithUsers := make([]*domain.ListingWithUser, 0, len(listings))
	for _, listing := range listings {
		user, exists := users[listing.UserID]
		if !exists {
			return nil, fmt.Errorf("user not found for listing %d", listing.ID)
		}

		listingWithUser := &domain.ListingWithUser{
			Listing: *listing,
			User:    *user,
		}
		listingsWithUsers = append(listingsWithUsers, listingWithUser)
	}

	return listingsWithUsers, nil
}

func (u *ListingUseCase) CreateListing(userID int, listingType string, price int) (*domain.Listing, error) {
	// Check if user exists
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Create listing
	return u.listingRepo.CreateListing(userID, listingType, price)
}
