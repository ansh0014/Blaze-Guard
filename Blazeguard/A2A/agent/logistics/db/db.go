package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type FireStation struct {
	ID       string
	Name     string
	Lat      float64
	Lng      float64
	Distance float64
}

type RouteInfo struct {
	StationID   string
	StationName string
	Distance    float64
	Duration    int64
	Geometry    string
	Steps       []RouteStep
}

type RouteStep struct {
	Instruction string
	Distance    float64
	Duration    int64
}

func InitDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),

	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("[Logistics] Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("[Logistics] Failed to ping database:", err)
	}

	fmt.Println("[Logistics] Connected to PostgreSQL with PostGIS")
}

func GetNearbyFireStations(location map[string]float64, radiusKm float64) []FireStation {
	lat := location["latitude"]
	lng := location["longitude"]

	query := `
        SELECT 
            id,
            name,
            ST_Y(location::geometry) as latitude,
            ST_X(location::geometry) as longitude,
            ST_Distance(
                location::geography,
                ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
            ) / 1000 as distance_km
        FROM fire_stations
        WHERE ST_DWithin(
            location::geography,
            ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
            $3 * 1000
        )
        ORDER BY distance_km ASC
        LIMIT 10;
    `

	rows, err := DB.Query(query, lng, lat, radiusKm)
	if err != nil {
		log.Printf("[Logistics] Database query error: %v", err)
		return []FireStation{}
	}
	defer rows.Close()

	var stations []FireStation
	for rows.Next() {
		var station FireStation
		err := rows.Scan(
			&station.ID,
			&station.Name,
			&station.Lat,
			&station.Lng,
			&station.Distance,
		)
		if err != nil {
			log.Printf("[Logistics] Row scan error: %v", err)
			continue
		}
		stations = append(stations, station)
	}

	fmt.Printf("[Logistics] Found %d fire stations within %.1fkm from database\n",
		len(stations), radiusKm)

	return stations
}

