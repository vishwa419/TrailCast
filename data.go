package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type WeatherDay struct {
	Date          string  `json:"date"`
	FormattedDate string  `json:"formattedDate"`
	MaxTemp       float64 `json:"maxTemp"`
	MinTemp       float64 `json:"minTemp"`
	Precipitation float64 `json:"precipitation"`
	WeatherCode   int     `json:"weatherCode"`
	Icon          string  `json:"icon"`
	IsWeekend     bool    `json:"isWeekend"`
	IsToday       bool    `json:"isToday"`
	RecommendDay  bool    `json:"recommendDay"` // <-- New field
}

type Park struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"imageUrl"`
	Description string  `json:"description"`
	License     string  `json:"license"`
	Attribution string  `json:"attribution"`
}

type Trail struct {
	Name         string  `json:"name"`
	Length       float64 `json:"length"`
	Difficulty   string  `json:"difficulty"`
	Description  string  `json:"description"`
	Elevation    float64 `json:"elevation"`
	EstTime      string  `json:"estTime"`
	ImageURL     string  `json:"imageUrl"`
	License      string  `json:"license"`
	Attribution  string  `json:"attribution"`
	StartLat     float64 `json:"startLat"`
	StartLng     float64 `json:"startLng"`
	RecommendDay int     `json:"recommendDay"`
}

func connectToDB() (*sql.DB, error) {
	connStr := "user=postgres password=hi dbname=parksdb sslmode=disable"
	return sql.Open("postgres", connStr)
}

func insertPark(db *sql.DB, park Park) (int, error) {
	var parkID int
	err := db.QueryRow(
		`INSERT INTO parks(name, location, latitude, longitude, image_url, description, license, attribution)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		park.Name, park.Location, park.Latitude, park.Longitude, park.ImageURL, park.Description, park.License, park.Attribution,
	).Scan(&parkID)

	if err != nil {
		return 0, err
	}
	return parkID, nil
}

func insertTrail(db *sql.DB, trail Trail, parkID int) error {
	_, err := db.Exec(
		`INSERT INTO trails(park_id, name, length, difficulty, description, elevation, estimated_time, image_url, license, attribution, start_latitude, start_longitude, recommend_day)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		parkID, trail.Name, trail.Length, trail.Difficulty, trail.Description, trail.Elevation, trail.EstTime, trail.ImageURL, trail.License, trail.Attribution, trail.StartLat, trail.StartLng, trail.RecommendDay,
	)
	return err
}

func fetchWeatherData(latitude, longitude float64) ([]WeatherDay, error) {
	// Simulate fetching data from an API
	return []WeatherDay{
		{
			Date:          "2025-05-01",
			FormattedDate: "Monday, May 1",
			MaxTemp:       22.5,
			MinTemp:       10.0,
			Precipitation: 0.0,
			WeatherCode:   0,
			Icon:          "sun",
			IsWeekend:     false,
			IsToday:       false,
		},
	}, nil
}

func getTrailsForPark(parkName string, lat, lng float64) ([]Trail, error) {
	rand.Seed(time.Now().UnixNano())

	// Sample trail data generation logic, similar to your current setup.
	trails := []Trail{
		{
			Name:         "Mountain Peak Trail",
			Length:       15.0,
			Difficulty:   "Moderate",
			Description:  "A scenic trail leading to the peak with amazing views.",
			Elevation:    500,
			EstTime:      "6-8 hours",
			ImageURL:     "https://example.com/image.jpg",
			License:      "CC BY-SA 4.0",
			Attribution:  "John Doe",
			StartLat:     lat,
			StartLng:     lng,
			RecommendDay: rand.Intn(7),
		},
	}

	return trails, nil
}

func main() {
	// Set up PostgreSQL connection
	db, err := connectToDB()
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	// Define parks
	parks := []Park{
		{
			Name:        "Yosemite National Park",
			Location:    "California",
			Latitude:    37.8651,
			Longitude:   -119.5383,
			ImageURL:    "https://example.com/yosemite.jpg",
			Description: "Known for waterfalls and stunning cliffs.",
			License:     "CC BY-SA 3.0",
			Attribution: "Wikipedia Contributors",
		},
	}

	// Insert parks and trails
	for _, park := range parks {
		// Insert park
		parkID, err := insertPark(db, park)
		if err != nil {
			log.Fatalf("Failed to insert park %s: %v", park.Name, err)
		}

		// Fetch trails for park
		trails, err := getTrailsForPark(park.Name, park.Latitude, park.Longitude)
		if err != nil {
			log.Fatalf("Failed to fetch trails for park %s: %v", park.Name, err)
		}

		// Insert trails into database
		for _, trail := range trails {
			err := insertTrail(db, trail, parkID)
			if err != nil {
				log.Fatalf("Failed to insert trail %s for park %s: %v", trail.Name, park.Name, err)
			}
		}
	}

	log.Println("Data insertion completed successfully!")
}
