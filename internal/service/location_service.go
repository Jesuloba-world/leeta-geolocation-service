package service

import (
	"log"
	"math"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/pkg/geospatial"
)

type LocationService struct {
	repo domain.LocationRepository
}

func NewLocationService(repo domain.LocationRepository) domain.LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) CreateLocation(name string, latitude, longitude float64) (*domain.Location, error) {
	log.Printf("Creating location: %s at (%.6f, %.6f)", name, latitude, longitude)

	location, err := domain.NewLocation(name, latitude, longitude)
	if err != nil {
		log.Printf("Failed to create location %s: %v", name, err)
		return nil, err
	}

	existing, _ := s.repo.FindByName(name)
	if existing != nil {
		log.Printf("Location %s already exists", name)
		return nil, domain.ErrLocationExists
	}

	err = s.repo.Save(location)
	if err != nil {
		log.Printf("Failed to save location %s: %v", name, err)
		return nil, err
	}

	log.Printf("Successfully created location: %s", name)
	return location, nil
}

func (s *LocationService) GetLocation(name string) (*domain.Location, error) {
	return s.repo.FindByName(name)
}

func (s *LocationService) GetAllLocations() ([]*domain.Location, error) {
	return s.repo.FindAll()
}

func (s *LocationService) DeleteLocation(name string) error {
	log.Printf("Deleting location: %s", name)
	err := s.repo.Delete(name)
	if err != nil {
		log.Printf("Failed to delete location %s: %v", name, err)
		return err
	}
	log.Printf("Successfully deleted location: %s", name)
	return nil
}

func (s *LocationService) FindNearest(latitude, longitude float64) (*domain.Location, float64, error) {
	locations, err := s.repo.FindAll()
	if err != nil {
		return nil, 0, err
	}

	if len(locations) == 0 {
		return nil, 0, domain.ErrLocationNotFound
	}

	var nearest *domain.Location
	minDistance := math.MaxFloat64

	for _, location := range locations {
		distance := geospatial.HaversineDistance(
			geospatial.Coordinate{Latitude: latitude, Longitude: longitude},
			geospatial.Coordinate{Latitude: location.Latitude, Longitude: location.Longitude},
		)

		if distance < minDistance {
			minDistance = distance
			nearest = location
		}
	}

	return nearest, minDistance, nil
}
