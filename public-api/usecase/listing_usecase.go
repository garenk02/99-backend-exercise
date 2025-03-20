package usecase

import (
	"fmt"
	"public-api/domain"
	"sync"
)

type ListingUseCase struct {
	listingRepo domain.ListingRepository
	userRepo    domain.UserRepository
	userCache   sync.Map // For caching users
}

func NewListingUseCase(listingRepo domain.ListingRepository, userRepo domain.UserRepository) *ListingUseCase {
	return &ListingUseCase{
		listingRepo: listingRepo,
		userRepo:    userRepo,
		userCache:   sync.Map{}, // Initialize the cache
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

	// Concurrent user fetching with caching
	users := make(map[int]*domain.User)
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect the users map

	errChan := make(chan error, len(userIDs)) // Buffered channel for errors

	for userID := range userIDs {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// Check cache first
			if cachedUser, ok := u.userCache.Load(userID); ok {
				mu.Lock()
				users[userID] = cachedUser.(*domain.User)
				mu.Unlock()
				return
			}

			user, err := u.userRepo.GetUserByID(userID)
			if err != nil {
				errChan <- fmt.Errorf("error fetching user data for listing: %w", err) // Send error to channel
				return
			}

			// Store in cache
			u.userCache.Store(userID, user)

			mu.Lock()
			users[userID] = user
			mu.Unlock()
		}(userID)
	}

	wg.Wait()
	close(errChan) // Close the error channel

	// Check for errors from goroutines
	for err := range errChan {
		return nil, err // Return the first error encountered
	}

	// Combine listings with user data using the cached users
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
