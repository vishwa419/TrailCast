package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Popular hiking parks with their coordinates

// HomeController handles the homepage routes
type HomeController struct {
	templates *template.Template
}

// NewHomeController creates a new home controller
func NewHomeController(templates *template.Template) *HomeController {
	return &HomeController{templates: templates}
}

// Index renders the home page
func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Parks []Park
		Title string
	}{
		Title: "Hiking Weather Forecast",
		Parks: make([]Park, 0, len(PopularParks)),
	}

	// Add parks to the slice
	for _, park := range PopularParks {
		data.Parks = append(data.Parks, park)
	}

	err := c.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// WeatherController handles weather related routes
type WeatherController struct {
	templates *template.Template
}

// NewWeatherController creates a new weather controller
func NewWeatherController(templates *template.Template) *WeatherController {
	return &WeatherController{
		templates: templates,
	}
}

// GetWeather renders the weather page
func (c *WeatherController) GetWeather(w http.ResponseWriter, r *http.Request) {
	parkName := r.URL.Query().Get("park")
	var park Park
	var ok bool

	if parkName != "" {
		park, ok = PopularParks[parkName]
	}

	// If park is not found or not specified, handle custom coordinates
	if !ok {
		latStr := r.URL.Query().Get("lat")
		lngStr := r.URL.Query().Get("lng")
		customName := r.URL.Query().Get("name")

		if latStr != "" && lngStr != "" {
			lat, errLat := strconv.ParseFloat(latStr, 64)
			lng, errLng := strconv.ParseFloat(lngStr, 64)

			if errLat == nil && errLng == nil {
				park = Park{
					Name:      customName,
					Latitude:  lat,
					Longitude: lng,
				}
			} else {
				http.Error(w, "Invalid coordinates", http.StatusBadRequest)
				return
			}
		} else {
			// Redirect to home if no park or coordinates are provided
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Fetch weather data
	weatherData, err := fetchWeatherData(park.Latitude, park.Longitude)
	if err != nil {
		http.Error(w, "Failed to fetch weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Choose the best weather day (simplified logic: high temp, low precipitation)
	bestIndex := -1
	bestScore := -1.0

	for i, day := range weatherData {
		score := (100 - day.Precipitation) + day.MaxTemp // You can adjust this formula
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	if bestIndex >= 0 {
		weatherData[bestIndex].RecommendDay = true
	}

	// Fetch trail data
	trails := getTrailsForPark(park.Name, park.Latitude, park.Longitude)

	// Create page data
	data := struct {
		Park         Park
		WeatherDays  []WeatherDay
		Trails       []Trail
		Title        string
		GeneratedAt  string
		OSMLatitude  float64
		OSMLongitude float64
	}{
		Park:         park,
		WeatherDays:  weatherData,
		Trails:       trails,
		Title:        fmt.Sprintf("Weather for %s", park.Name),
		GeneratedAt:  time.Now().Format("January 2, 2006 at 3:04 PM"),
		OSMLatitude:  park.Latitude,
		OSMLongitude: park.Longitude,
	}

	err = c.templates.ExecuteTemplate(w, "weather.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetWeatherAPI handles API requests for weather data
func (c *WeatherController) GetWeatherAPI(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Latitude and longitude are required", http.StatusBadRequest)
		return
	}

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)

	if errLat != nil || errLng != nil {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	// Fetch weather data
	weatherData, err := fetchWeatherData(lat, lng)
	if err != nil {
		http.Error(w, "Failed to fetch weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	json.NewEncoder(w).Encode(weatherData)
}

// TrailController handles trail related routes
type TrailController struct {
	templates *template.Template
}

// NewTrailController creates a new trail controller
func NewTrailController(templates *template.Template) *TrailController {
	return &TrailController{
		templates: templates,
	}
}

// GetTrailsAPI handles API requests for trail data
func (c *TrailController) GetTrailsAPI(w http.ResponseWriter, r *http.Request) {
	parkName := r.URL.Query().Get("park")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	var lat, lng float64
	var err1, err2 error

	if latStr != "" && lngStr != "" {
		lat, err1 = strconv.ParseFloat(latStr, 64)
		lng, err2 = strconv.ParseFloat(lngStr, 64)
	} else if parkName != "" {
		if park, ok := PopularParks[parkName]; ok {
			lat = park.Latitude
			lng = park.Longitude
		} else {
			http.Error(w, "Park not found", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Park name or coordinates are required", http.StatusBadRequest)
		return
	}

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	// Get trails
	trails := getTrailsForPark(parkName, lat, lng)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	json.NewEncoder(w).Encode(trails)
}
