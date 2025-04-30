package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// Application configuration
type Config struct {
	Port string
}

// Initialize application configuration
func initConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return Config{
		Port: port,
	}
}

func main() {
	config := initConfig()

	// Create templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Initialize router and controllers
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static/css"))
	mux.Handle("/static/css/", http.StripPrefix("/static/css/", fs))
	
	// Register controllers/handlers
	homeController := NewHomeController(templates)
	weatherController := NewWeatherController(templates)
	trailController := NewTrailController(templates)

	// Register routes
	mux.HandleFunc("/", homeController.Index)
	mux.HandleFunc("/weather", weatherController.GetWeather)
	mux.HandleFunc("/api/weather", weatherController.GetWeatherAPI)
	mux.HandleFunc("/api/trails", trailController.GetTrailsAPI)

	// Start server
	serverAddr := fmt.Sprintf(":%s", config.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting server on http://localhost:%s\n", config.Port)
	log.Fatal(server.ListenAndServe())
}
