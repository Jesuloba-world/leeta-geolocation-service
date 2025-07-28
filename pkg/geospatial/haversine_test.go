package geospatial

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		p1       Coordinate
		p2       Coordinate
		expected float64
		delta    float64
	}{
		{
			name: "New York to Los Angeles",
			p1: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			p2: Coordinate{
				Latitude:  34.0522,
				Longitude: -118.2437,
			},
			expected: 3935.9,
			delta:    0.5,
		},
		{
			name: "London to Paris",
			p1: Coordinate{
				Latitude:  51.5074,
				Longitude: -0.1278,
			},
			p2: Coordinate{
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			expected: 343.6,
			delta:    0.5,
		},
		{
			name: "Same point",
			p1: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			p2: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expected: 0,
			delta:    0.001,
		},
		{
			name: "Antipodal points",
			p1: Coordinate{
				Latitude:  0,
				Longitude: 0,
			},
			p2: Coordinate{
				Latitude:  0,
				Longitude: 180,
			},
			expected: 20015.1,
			delta:    0.5,
		},
		{
			name: "Short distance",
			p1: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			p2: Coordinate{
				Latitude:  40.7589,
				Longitude: -73.9851,
			},
			expected: 5.4,
			delta:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := HaversineDistance(tt.p1, tt.p2)
			if math.Abs(distance-tt.expected) > tt.delta {
				t.Errorf("HaversineDistance() = %v, want %v (±%v)", distance, tt.expected, tt.delta)
			}
		})
	}
}

func TestToRadians(t *testing.T) {
	tests := []struct {
		name     string
		degrees  float64
		expected float64
		delta    float64
	}{
		{
			name:     "Zero degrees",
			degrees:  0,
			expected: 0,
			delta:    0.001,
		},
		{
			name:     "90 degrees",
			degrees:  90,
			expected: math.Pi / 2,
			delta:    0.001,
		},
		{
			name:     "180 degrees",
			degrees:  180,
			expected: math.Pi,
			delta:    0.001,
		},
		{
			name:     "360 degrees",
			degrees:  360,
			expected: 2 * math.Pi,
			delta:    0.001,
		},
		{
			name:     "Negative degrees",
			degrees:  -90,
			expected: -math.Pi / 2,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			radians := toRadians(tt.degrees)
			if math.Abs(radians-tt.expected) > tt.delta {
				t.Errorf("toRadians(%v) = %v, want %v (±%v)", tt.degrees, radians, tt.expected, tt.delta)
			}
		})
	}
}

func TestKmToMiles(t *testing.T) {
	tests := []struct {
		name     string
		km       float64
		expected float64
		delta    float64
	}{
		{
			name:     "Zero km",
			km:       0,
			expected: 0,
			delta:    0.001,
		},
		{
			name:     "1 km",
			km:       1,
			expected: 0.621371,
			delta:    0.001,
		},
		{
			name:     "100 km",
			km:       100,
			expected: 62.1371,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			miles := KmToMiles(tt.km)
			if math.Abs(miles-tt.expected) > tt.delta {
				t.Errorf("KmToMiles(%v) = %v, want %v (±%v)", tt.km, miles, tt.expected, tt.delta)
			}
		})
	}
}

func TestKmToNauticalMiles(t *testing.T) {
	tests := []struct {
		name     string
		km       float64
		expected float64
		delta    float64
	}{
		{
			name:     "Zero km",
			km:       0,
			expected: 0,
			delta:    0.001,
		},
		{
			name:     "1 km",
			km:       1,
			expected: 0.539957,
			delta:    0.001,
		},
		{
			name:     "100 km",
			km:       100,
			expected: 53.9957,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nauticalMiles := KmToNauticalMiles(tt.km)
			if math.Abs(nauticalMiles-tt.expected) > tt.delta {
				t.Errorf("KmToNauticalMiles(%v) = %v, want %v (±%v)", tt.km, nauticalMiles, tt.expected, tt.delta)
			}
		})
	}
}

func TestMilesToKm(t *testing.T) {
	tests := []struct {
		name     string
		miles    float64
		expected float64
		delta    float64
	}{
		{
			name:     "Zero miles",
			miles:    0,
			expected: 0,
			delta:    0.001,
		},
		{
			name:     "1 mile",
			miles:    1,
			expected: 1.609344,
			delta:    0.001,
		},
		{
			name:     "100 miles",
			miles:    100,
			expected: 160.9344,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := MilesToKm(tt.miles)
			if math.Abs(km-tt.expected) > tt.delta {
				t.Errorf("MilesToKm(%v) = %v, want %v (±%v)", tt.miles, km, tt.expected, tt.delta)
			}
		})
	}
}

func TestNauticalMilesToKm(t *testing.T) {
	tests := []struct {
		name          string
		nauticalMiles float64
		expected      float64
		delta         float64
	}{
		{
			name:          "Zero nautical miles",
			nauticalMiles: 0,
			expected:      0,
			delta:         0.001,
		},
		{
			name:          "1 nautical mile",
			nauticalMiles: 1,
			expected:      1.852,
			delta:         0.001,
		},
		{
			name:          "100 nautical miles",
			nauticalMiles: 100,
			expected:      185.2,
			delta:         0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := NauticalMilesToKm(tt.nauticalMiles)
			if math.Abs(km-tt.expected) > tt.delta {
				t.Errorf("NauticalMilesToKm(%v) = %v, want %v (±%v)", tt.nauticalMiles, km, tt.expected, tt.delta)
			}
		})
	}
}

func TestHaversineDistanceMiles(t *testing.T) {
	p1 := Coordinate{Latitude: 40.7128, Longitude: -74.0060}
	p2 := Coordinate{Latitude: 34.0522, Longitude: -118.2437}

	distanceKm := HaversineDistance(p1, p2)
	distanceMiles := HaversineDistanceMiles(p1, p2)
	expectedMiles := KmToMiles(distanceKm)

	if math.Abs(distanceMiles-expectedMiles) > 0.001 {
		t.Errorf("HaversineDistanceMiles() = %v, want %v", distanceMiles, expectedMiles)
	}
}

func TestHaversineDistanceNauticalMiles(t *testing.T) {
	p1 := Coordinate{Latitude: 40.7128, Longitude: -74.0060}
	p2 := Coordinate{Latitude: 34.0522, Longitude: -118.2437}

	distanceKm := HaversineDistance(p1, p2)
	distanceNauticalMiles := HaversineDistanceNauticalMiles(p1, p2)
	expectedNauticalMiles := KmToNauticalMiles(distanceKm)

	if math.Abs(distanceNauticalMiles-expectedNauticalMiles) > 0.001 {
		t.Errorf("HaversineDistanceNauticalMiles() = %v, want %v", distanceNauticalMiles, expectedNauticalMiles)
	}
}
