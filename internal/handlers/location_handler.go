package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/pkg/errors"
)

type LocationHandler struct {
	service domain.LocationService
}

func NewLocationHandler(service domain.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// CreateLocation handles POST /locations requests
func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.RespondWithError(w, errors.BadRequest("Method not allowed"))
		return
	}

	var location domain.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		errors.RespondWithError(w, errors.BadRequest("Invalid request body"))
		return
	}

	createdLocation, err := h.service.CreateLocation(location.Name, location.Latitude, location.Longitude)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			errors.RespondWithError(w, errors.Conflict(err.Error()))
			return
		}
		errors.RespondWithError(w, errors.BadRequest(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLocation)
}

// GetAllLocations handles GET /locations requests
func (h *LocationHandler) GetAllLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.RespondWithError(w, errors.BadRequest("Method not allowed"))
		return
	}

	locations, err := h.service.GetAllLocations()
	if err != nil {
		errors.RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

// DeleteLocation handles DELETE /locations/{name} requests
func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errors.RespondWithError(w, errors.BadRequest("Method not allowed"))
		return
	}

	// Extract name from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		errors.RespondWithError(w, errors.BadRequest("Invalid URL path"))
		return
	}

	name := pathParts[2] // /locations/{name}

	if err := h.service.DeleteLocation(name); err != nil {
		if strings.Contains(err.Error(), "not found") {
			errors.RespondWithError(w, errors.NotFound(err.Error()))
			return
		}
		errors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FindNearest handles GET /nearest?lat=LAT&lng=LNG requests
func (h *LocationHandler) FindNearest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.RespondWithError(w, errors.BadRequest("Method not allowed"))
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	latStr := query.Get("lat")
	lngStr := query.Get("lng")

	if latStr == "" || lngStr == "" {
		errors.RespondWithError(w, errors.BadRequest("Missing lat or lng parameters"))
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		errors.RespondWithError(w, errors.BadRequest("Invalid latitude value"))
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		errors.RespondWithError(w, errors.BadRequest("Invalid longitude value"))
		return
	}

	nearestLocation, distance, err := h.service.FindNearest(lat, lng)
	if err != nil {
		if strings.Contains(err.Error(), "no locations") {
			errors.RespondWithError(w, errors.NotFound(err.Error()))
			return
		}
		errors.RespondWithError(w, err)
		return
	}

	response := struct {
		Location *domain.Location `json:"location"`
		Distance float64         `json:"distance_km"`
	}{
		Location: nearestLocation,
		Distance: distance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}