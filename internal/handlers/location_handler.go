package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/internal/dto"
)

// LocationRequest represents the request body for creating a location
type LocationRequest struct {
	Body dto.LocationRequest `json:"body"`
}

// LocationResponse represents a location response
type LocationResponse struct {
	Body dto.LocationResponse `json:"body"`
}

// LocationListResponse represents a list of locations
type LocationListResponse struct {
	Body dto.LocationListResponse `json:"body"`
}

// NearestLocationRequest represents the query parameters for finding nearest location
type NearestLocationRequest struct {
	Lat float64 `query:"lat" required:"true" minimum:"-90" maximum:"90" doc:"Latitude coordinate"`
	Lng float64 `query:"lng" required:"true" minimum:"-180" maximum:"180" doc:"Longitude coordinate"`
}

// NearestLocationResponse represents the nearest location response
type NearestLocationResponse struct {
	Body dto.NearestLocationResponse `json:"body"`
}

// DeleteLocationRequest represents the path parameter for deleting a location
type DeleteLocationRequest struct {
	Name string `path:"name" required:"true" doc:"Name of the location to delete"`
}

// HealthResponse represents the health check response
// LocationHandler wraps the location service for API operations
type LocationHandler struct {
	service domain.LocationService
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(service domain.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// RegisterRoutes registers all location routes with the Huma API
func (h *LocationHandler) RegisterRoutes(api huma.API) {
	// Create location endpoint
	huma.Register(api, huma.Operation{
		OperationID:   "create-location",
		Method:        http.MethodPost,
		Path:          "/locations",
		Summary:       "Create Location",
		Description:   "Register a new geolocated station with latitude and longitude coordinates",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusCreated,
	}, h.CreateLocation)

	// Get all locations endpoint
	huma.Register(api, huma.Operation{
		OperationID: "get-locations",
		Method:      http.MethodGet,
		Path:        "/locations",
		Summary:     "Get All Locations",
		Description: "Retrieve all registered locations",
		Tags:        []string{"Locations"},
	}, h.GetAllLocations)

	// Delete location endpoint
	huma.Register(api, huma.Operation{
		OperationID:   "delete-location",
		Method:        http.MethodDelete,
		Path:          "/locations/{name}",
		Summary:       "Delete Location",
		Description:   "Delete a location by its unique name",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusNoContent,
	}, h.DeleteLocation)

	// Find nearest location endpoint
	huma.Register(api, huma.Operation{
		OperationID: "find-nearest",
		Method:      http.MethodGet,
		Path:        "/nearest",
		Summary:     "Find Nearest Location",
		Description: "Find the closest registered location to the given coordinates",
		Tags:        []string{"Locations"},
	}, h.FindNearest)
}

// CreateLocation handles POST /locations requests
func (h *LocationHandler) CreateLocation(ctx context.Context, input *LocationRequest) (*LocationResponse, error) {
	createdLocation, err := h.service.CreateLocation(input.Body.Name, input.Body.Latitude, input.Body.Longitude)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, huma.Error409Conflict("Location with this name already exists")
		}
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &LocationResponse{
		Body: dto.FromDomain(createdLocation),
	}, nil
}

// GetAllLocations handles GET /locations requests
func (h *LocationHandler) GetAllLocations(ctx context.Context, input *struct{}) (*LocationListResponse, error) {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve locations")
	}

	return &LocationListResponse{
		Body: dto.FromDomainList(locations),
	}, nil
}

// DeleteLocation handles DELETE /locations/{name} requests
func (h *LocationHandler) DeleteLocation(ctx context.Context, input *DeleteLocationRequest) (*struct{}, error) {
	err := h.service.DeleteLocation(input.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, huma.Error404NotFound("Location not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete location")
	}

	return &struct{}{}, nil
}

// FindNearest handles GET /nearest requests
func (h *LocationHandler) FindNearest(ctx context.Context, input *NearestLocationRequest) (*NearestLocationResponse, error) {
	location, distance, err := h.service.FindNearest(input.Lat, input.Lng)
	if err != nil {
		if strings.Contains(err.Error(), "no locations") {
			return nil, huma.Error404NotFound("No locations found")
		}
		return nil, huma.Error500InternalServerError("Failed to find nearest location")
	}

	return &NearestLocationResponse{
		Body: dto.FromDomainWithDistance(location, distance),
	}, nil
}