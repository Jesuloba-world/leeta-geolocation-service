package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jesuloba-world/leeta-task/pkg/validator"
)

type Location struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=1"`
	Latitude  float64   `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64   `json:"longitude" validate:"required,min=-180,max=180"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrEmptyName        = errors.New("location name cannot be empty")
	ErrInvalidLatitude  = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")
	ErrLocationNotFound = errors.New("location not found")
	ErrLocationExists   = errors.New("location already exists")
)

func NewLocation(name string, latitude, longitude float64) (*Location, error) {
	location := &Location{
		Name:      strings.TrimSpace(name),
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
	}

	if err := location.Validate(); err != nil {
		return nil, err
	}

	return location, nil
}

func (l *Location) Validate() error {
	return validator.ValidateStruct(l)
}

func (l *Location) String() string {
	return l.Name + " (" + formatCoordinate(l.Latitude) + ", " + formatCoordinate(l.Longitude) + ")"
}

func formatCoordinate(coord float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", coord), "0"), ".")
}

type LocationRepository interface {
	Save(location *Location) error
	FindByName(name string) (*Location, error)
	FindAll() ([]*Location, error)
	Delete(name string) error
}

type LocationService interface {
	CreateLocation(name string, latitude, longitude float64) (*Location, error)
	GetLocation(name string) (*Location, error)
	GetAllLocations() ([]*Location, error)
	DeleteLocation(name string) error
	FindNearest(latitude, longitude float64) (*Location, float64, error)
}
