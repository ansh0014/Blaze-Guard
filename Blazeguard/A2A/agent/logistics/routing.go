package main

import (
	"blazeguard/agent/logistics/config"
	"blazeguard/agent/logistics/db"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

type routeResult struct {
	route database.RouteInfo
	err   error
}

func calculateBestRouteParallel(fireLocation map[string]float64, fireStations []database.FireStation) (database.RouteInfo, error) {
	if len(fireStations) == 0 {
		return database.RouteInfo{}, fmt.Errorf("no fire stations available")
	}

	nearbyStations := filterStationsByHaversine(fireLocation, fireStations, 50.0)
	if len(nearbyStations) == 0 {
		return database.RouteInfo{}, fmt.Errorf("no stations within range")
	}

	fmt.Printf("[Logistics] Calculating routes for %d stations in parallel\n", len(nearbyStations))

	resultChan := make(chan routeResult, len(nearbyStations))
	var wg sync.WaitGroup

	for _, station := range nearbyStations {
		wg.Add(1)
		go func(s database.FireStation) {
			defer wg.Done()

			route, err := calculateRouteMapboxWithRetry(fireLocation, s)
			resultChan <- routeResult{route: route, err: err}
		}(station)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var bestRoute database.RouteInfo
	minDuration := int64(999999)
	routesCalculated := 0
	routesFailed := 0

	for result := range resultChan {
		if result.err != nil {
			routesFailed++
			continue
		}

		routesCalculated++
		if result.route.Duration < minDuration {
			minDuration = result.route.Duration
			bestRoute = result.route
		}
	}

	fmt.Printf("[Logistics] Routes calculated: %d success, %d failed\n", routesCalculated, routesFailed)

	if bestRoute.StationID == "" {
		return database.RouteInfo{}, fmt.Errorf("all route calculations failed")
	}

	return bestRoute, nil
}

func filterStationsByHaversine(fireLocation map[string]float64, stations []database.FireStation, maxRadiusKm float64) []database.FireStation {
	fireLat := fireLocation["latitude"]
	fireLng := fireLocation["longitude"]

	filtered := []database.FireStation{}
	for _, station := range stations {
		distance := haversineDistance(fireLat, fireLng, station.Lat, station.Lng)
		if distance <= maxRadiusKm {
			station.Distance = distance
			filtered = append(filtered, station)
		}
	}

	fmt.Printf("[Logistics] Haversine filter: %d/%d stations within %.1fkm\n",
		len(filtered), len(stations), maxRadiusKm)

	return filtered
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return config.EARTH_RADIUS_KM * c
}

func calculateRouteMapboxWithRetry(fireLocation map[string]float64, station database.FireStation) (database.RouteInfo, error) {
	var lastErr error

	for attempt := 1; attempt <= config.MAX_RETRIES; attempt++ {
		<-config.MapboxRateLimiter

		route, err := calculateRouteMapbox(fireLocation, station)

		config.MapboxRateLimiter <- struct{}{}

		if err == nil {
			return route, nil
		}

		lastErr = err
		if attempt < config.MAX_RETRIES {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	return database.RouteInfo{}, fmt.Errorf("failed after %d retries: %v", config.MAX_RETRIES, lastErr)
}

func calculateRouteMapbox(fireLocation map[string]float64, station database.FireStation) (database.RouteInfo, error) {
	apiKey := os.Getenv("MAPBOX_API_KEY")
	fireLat := fireLocation["latitude"]
	fireLng := fireLocation["longitude"]

	url := fmt.Sprintf(
		"https://api.mapbox.com/directions/v5/mapbox/driving-traffic/%f,%f;%f,%f?geometries=geojson&steps=true&access_token=%s",
		station.Lng, station.Lat, fireLng, fireLat, apiKey,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return database.RouteInfo{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.TIMEOUT_SECONDS*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := config.HTTPClient.Do(req)
	if err != nil {
		return database.RouteInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return database.RouteInfo{}, fmt.Errorf("mapbox API error: status %d", resp.StatusCode)
	}

	var result MapboxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return database.RouteInfo{}, err
	}

	if len(result.Routes) == 0 {
		return database.RouteInfo{}, fmt.Errorf("no routes found")
	}

	route := result.Routes[0]
	steps := []database.RouteStep{}

	if len(route.Legs) > 0 {
		for _, step := range route.Legs[0].Steps {
			steps = append(steps, database.RouteStep{
				Instruction: step.Maneuver.Instruction,
				Distance:    step.Distance,
				Duration:    int64(step.Duration),
			})
		}
	}

	geometryJSON, _ := json.Marshal(route.Geometry)

	return database.RouteInfo{
		StationID:   station.ID,
		StationName: station.Name,
		Distance:    route.Distance,
		Duration:    int64(route.Duration),
		Geometry:    string(geometryJSON),
		Steps:       steps,
	}, nil
}

type MapboxResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Legs []struct {
			Steps []struct {
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
				Maneuver struct {
					Instruction string `json:"instruction"`
				} `json:"maneuver"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
}
