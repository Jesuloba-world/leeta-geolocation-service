package geospatial

import (
	"math"
)

// EarthRadiusKm is the radius of the Earth in kilometers
const EarthRadiusKm = 6371.0

// Coordinate represents a geographic point with latitude and longitude
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// toRadians converts degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// HaversineDistance calculates the distance between two coordinates using the Haversine formula
// Returns distance in kilometers
func HaversineDistance(p1, p2 Coordinate) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1 := toRadians(p1.Latitude)
	lon1 := toRadians(p1.Longitude)
	lat2 := toRadians(p2.Latitude)
	lon2 := toRadians(p2.Longitude)

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := EarthRadiusKm * c

	return distance
}

// Conversion constants
const (
	KmToMilesRatio        = 0.621371
	KmToNauticalMilesRatio = 0.539957
	MilesToKmRatio        = 1.609344
	NauticalMilesToKmRatio = 1.852
)

// KmToMiles converts kilometers to miles
func KmToMiles(km float64) float64 {
	return km * KmToMilesRatio
}

// KmToNauticalMiles converts kilometers to nautical miles
func KmToNauticalMiles(km float64) float64 {
	return km * KmToNauticalMilesRatio
}

// MilesToKm converts miles to kilometers
func MilesToKm(miles float64) float64 {
	return miles * MilesToKmRatio
}

// NauticalMilesToKm converts nautical miles to kilometers
func NauticalMilesToKm(nauticalMiles float64) float64 {
	return nauticalMiles * NauticalMilesToKmRatio
}

// HaversineDistanceMiles calculates distance in miles
func HaversineDistanceMiles(p1, p2 Coordinate) float64 {
	return KmToMiles(HaversineDistance(p1, p2))
}

// HaversineDistanceNauticalMiles calculates distance in nautical miles
func HaversineDistanceNauticalMiles(p1, p2 Coordinate) float64 {
	return KmToNauticalMiles(HaversineDistance(p1, p2))
}