// handlers/listing_handler.go
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"public-api/domain"
	"strconv"
)

type ListingHandler struct {
	listingUseCase domain.ListingUseCase
}

func NewListingHandler(listingUseCase domain.ListingUseCase) *ListingHandler {
	return &ListingHandler{
		listingUseCase: listingUseCase,
	}
}

func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	
	// Parse page_num
	pageNum := 1
	if pageNumStr := query.Get("page_num"); pageNumStr != "" {
		if num, err := strconv.Atoi(pageNumStr); err == nil && num > 0 {
			pageNum = num
		} else if err != nil {
			domain.RespondWithError(w, http.StatusBadRequest, "Invalid page_num parameter", err)
			return
		}
	}
	
	// Parse page_size
	pageSize := 10
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if size, err := strconv.Atoi(pageSizeStr); err == nil && size > 0 {
			pageSize = size
		} else if err != nil {
			domain.RespondWithError(w, http.StatusBadRequest, "Invalid page_size parameter", err)
			return
		}
	}
	
	// Parse user_id
	var userID *int
	if userIDStr := query.Get("user_id"); userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		} else {
			domain.RespondWithError(w, http.StatusBadRequest, "Invalid user_id parameter", err)
			return
		}
	}
	
	// Log request parameters
	slog.Info("Fetching listings", 
		"page_num", pageNum, 
		"page_size", pageSize, 
		"user_id", userID,
	)
	
	// Get listings
	listings, err := h.listingUseCase.GetListings(pageNum, pageSize, userID)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch listings", err)
		return
	}
	
	// Log success
	slog.Info("Listings fetched successfully", "count", len(listings))
	
	// Prepare response
	response := struct {
		Result   bool                    `json:"result"`
		Listings []*domain.ListingWithUser `json:"listings"`
	}{
		Result:   true,
		Listings: listings,
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var request struct {
		UserID      int    `json:"user_id"`
		ListingType string `json:"listing_type"`
		Price       int    `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		domain.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if request.UserID <= 0 {
		domain.RespondWithError(w, http.StatusBadRequest, "Invalid user_id", nil)
		return
	}
	if request.ListingType == "" {
		domain.RespondWithError(w, http.StatusBadRequest, "listing_type is required", nil)
		return
	}
	if request.Price <= 0 {
		domain.RespondWithError(w, http.StatusBadRequest, "Invalid price", nil)
		return
	}

	// Log request
	slog.Info("Creating listing", 
		"user_id", request.UserID, 
		"listing_type", request.ListingType, 
		"price", request.Price,
	)

	// Create listing
	listing, err := h.listingUseCase.CreateListing(request.UserID, request.ListingType, request.Price)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, "Failed to create listing", err)
		return
	}

	// Log success
	slog.Info("Listing created successfully", "listing_id", listing.ID)

	// Prepare response
	response := struct {
		Listing *domain.Listing `json:"listing"`
	}{
		Listing: listing,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}