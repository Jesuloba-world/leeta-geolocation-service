package memory

import (
	"fmt"
	"math"
	"sync"

	"github.com/jesuloba-world/leeta-task/internal/domain"
	"github.com/jesuloba-world/leeta-task/pkg/geospatial"
)

type InMemoryLocationRepository struct {
	mu        sync.RWMutex
	locations map[string]*domain.Location // key is name
	locationsById map[string]*domain.Location // key is ID
	nextID    int
}

func NewInMemoryLocationRepository() *InMemoryLocationRepository {
	return &InMemoryLocationRepository{
		locations: make(map[string]*domain.Location),
		locationsById: make(map[string]*domain.Location),
		nextID:    1,
	}
}

func (r *InMemoryLocationRepository) Save(location *domain.Location) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.locations[location.Name]; exists {
		return domain.ErrLocationExists
	}

	if location.ID == "" {
		location.ID = fmt.Sprintf("%d", r.nextID)
		r.nextID++
	}

	r.locations[location.Name] = location
	r.locationsById[location.ID] = location
	return nil
}

func (r *InMemoryLocationRepository) FindByName(name string) (*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	location, exists := r.locations[name]
	if !exists {
		return nil, domain.ErrLocationNotFound
	}

	return location, nil
}

func (r *InMemoryLocationRepository) FindAll() ([]*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	locations := make([]*domain.Location, 0, len(r.locations))
	for _, location := range r.locations {
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *InMemoryLocationRepository) Delete(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.locations[name]; !exists {
		return domain.ErrLocationNotFound
	}

	delete(r.locations, name)
	return nil
}

func (r *InMemoryLocationRepository) FindNearest(latitude, longitude float64) (*domain.Location, float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.locations) == 0 {
		return nil, 0, domain.ErrLocationNotFound
	}

	var nearest *domain.Location
	minDistance := math.MaxFloat64

	for _, location := range r.locations {
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

func (r *InMemoryLocationRepository) FindByID(id string) (*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	location, exists := r.locationsById[id]
	if !exists {
		return nil, domain.ErrLocationNotFound
	}

	return location, nil
}