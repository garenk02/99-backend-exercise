package repository

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"public-api/domain"
	"strconv"
	"strings"
)

type ListingRepository struct {
	baseURL string
}

func NewListingRepository(baseURL string) *ListingRepository {
	return &ListingRepository{
		baseURL: baseURL,
	}
}

func (r *ListingRepository) GetListings(pageNum, pageSize int, userID *int) ([]*domain.Listing, error) {
	// Build URL with query parameters
	reqURL := fmt.Sprintf("%s/listings?page_num=%d&page_size=%d", r.baseURL, pageNum, pageSize)
	if userID != nil {
		reqURL = fmt.Sprintf("%s&user_id=%d", reqURL, *userID)
	}

	slog.Debug("Fetching listings",
		"page_num", pageNum,
		"page_size", pageSize,
		"user_id", userID,
		"url", reqURL,
	)

	// Make HTTP request
	resp, err := http.Get(reqURL)
	if err != nil {
		slog.Error("Error making request to listing service", "error", err)
		return nil, fmt.Errorf("error making request to listing service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		slog.Error("Listing service returned non-200 status", "status", resp.StatusCode)
		return nil, fmt.Errorf("listing service returned non-200 status: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Result   bool `json:"result"`
		Listings []struct {
			ID          int    `json:"id"`
			UserID      int    `json:"user_id"`
			ListingType string `json:"listing_type"`
			Price       int    `json:"price"`
			CreatedAt   int64  `json:"created_at"`
			UpdatedAt   int64  `json:"updated_at"`
		} `json:"listings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding response from listing service", "error", err)
		return nil, fmt.Errorf("error decoding response from listing service: %w", err)
	}

	// Convert to domain model
	listings := make([]*domain.Listing, 0, len(response.Listings))
	for _, l := range response.Listings {
		listings = append(listings, &domain.Listing{
			ID:          l.ID,
			UserID:      l.UserID,
			ListingType: l.ListingType,
			Price:       l.Price,
			CreatedAt:   l.CreatedAt,
			UpdatedAt:   l.UpdatedAt,
		})
	}

	slog.Debug("Fetched listings successfully", "count", len(listings))
	return listings, nil
}

func (r *ListingRepository) CreateListing(userID int, listingType string, price int) (*domain.Listing, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("user_id", strconv.Itoa(userID))
	data.Set("listing_type", listingType)
	data.Set("price", strconv.Itoa(price))

	reqURL := fmt.Sprintf("%s/listings", r.baseURL)
	slog.Debug("Creating listing",
		"user_id", userID,
		"listing_type", listingType,
		"price", price,
		"url", reqURL,
	)

	// Make HTTP request
	resp, err := http.Post(
		reqURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		slog.Error("Error making request to listing service", "error", err)
		return nil, fmt.Errorf("error making request to listing service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		slog.Error("Listing service returned status", "status", resp.StatusCode)
		return nil, fmt.Errorf("listing service returned status: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Result  bool `json:"result"`
		Listing struct {
			ID          int    `json:"id"`
			UserID      int    `json:"user_id"`
			ListingType string `json:"listing_type"`
			Price       int    `json:"price"`
			CreatedAt   int64  `json:"created_at"`
			UpdatedAt   int64  `json:"updated_at"`
		} `json:"listing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding response from listing service", "error", err)
		return nil, fmt.Errorf("error decoding response from listing service: %w", err)
	}

	// Convert to domain model
	listing := &domain.Listing{
		ID:          response.Listing.ID,
		UserID:      response.Listing.UserID,
		ListingType: response.Listing.ListingType,
		Price:       response.Listing.Price,
		CreatedAt:   response.Listing.CreatedAt,
		UpdatedAt:   response.Listing.UpdatedAt,
	}

	slog.Debug("Listing created successfully", "listing_id", listing.ID)
	return listing, nil
}
