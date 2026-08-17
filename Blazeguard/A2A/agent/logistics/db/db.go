package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

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
	var connStr string

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			sslMode = "disable"
		}
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			sslMode,
		)
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("[Logistics] Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("[Logistics] Failed to ping database:", err)
	}

	fmt.Println("[Logistics] Connected to PostgreSQL with PostGIS")

	if err := runMigrations(); err != nil {
		log.Fatal("[Logistics] PostgreSQL migrations failed:", err)
	}
}

func runMigrations() error {
	fmt.Println("[Logistics] Running database migrations...")
	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := "migrations/" + file.Name()
		fmt.Printf("[Logistics] Applying migration: %s\n", file.Name())

		content, err := migrationFS.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file.Name(), err)
		}
	}

	fmt.Println("[Logistics] PostgreSQL migrations applied successfully")
	return nil
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

