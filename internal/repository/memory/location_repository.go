package memory

import (
	"fmt"
	"sync"
	"github.com/jesuloba-world/leeta-task/internal/domain"
)

type InMemoryLocationRepository struct {
	mu        sync.RWMutex
	locations map[string]*domain.Location
	nextID    int
}

func NewInMemoryLocationRepository() *InMemoryLocationRepository {
	return &InMemoryLocationRepository{
		locations: make(map[string]*domain.Location),
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