package dto

import (
	"time"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/pkg/validator"
)

type LocationRequest struct {
	Name      string  `json:"name" binding:"required" validate:"required,min=1"`
	Latitude  float64 `json:"latitude" binding:"required" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required" validate:"required,min=-180,max=180"`
}

type LocationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type LocationListResponse struct {
	Locations []LocationResponse `json:"locations"`
	Count     int                `json:"count"`
}

type NearestLocationResponse struct {
	Location LocationResponse `json:"location"`
	Distance float64          `json:"distance_km"`
}

func (req *LocationRequest) Validate() error {
	return validator.ValidateStruct(req)
}

func (req *LocationRequest) ToDomain() (*domain.Location, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return domain.NewLocation(req.Name, req.Latitude, req.Longitude)
}

func FromDomain(location *domain.Location) LocationResponse {
	return LocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		CreatedAt: location.CreatedAt,
	}
}

func FromDomainList(locations []*domain.Location) LocationListResponse {
	responses := make([]LocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = FromDomain(location)
	}

	return LocationListResponse{
		Locations: responses,
		Count:     len(responses),
	}
}

func FromDomainWithDistance(location *domain.Location, distance float64) NearestLocationResponse {
	return NearestLocationResponse{
		Location: FromDomain(location),
		Distance: distance,
	}
}
