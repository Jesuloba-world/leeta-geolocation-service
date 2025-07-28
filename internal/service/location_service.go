package service

import (
	"log"

	"github.com/jesuloba-world/leeta-task/internal/domain"
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

func (s *LocationService) GetLocationByID(id string) (*domain.Location, error) {
	return s.repo.FindByID(id)
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
	return s.repo.FindNearest(latitude, longitude)
}
