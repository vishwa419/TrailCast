package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Park represents a hiking park location
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

// WeatherDay represents a single day's weather forecast
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

// Trail represents a hiking trail
type Trail struct {
	Name         string  `json:"name"`
	Length       float64 `json:"length"` // in kilometers
	Difficulty   string  `json:"difficulty"`
	Description  string  `json:"description"`
	Elevation    float64 `json:"elevation"` // in meters
	EstTime      string  `json:"estTime"`
	ImageURL     string  `json:"imageUrl"`
	License      string  `json:"license"`
	Attribution  string  `json:"attribution"`
	StartLat     float64 `json:"startLat"`
	StartLng     float64 `json:"startLng"`
	RecommendDay int     `json:"recommendDay"` // Index of recommended day based on weather
}

// OpenMeteoResponse represents the response from Open-Meteo API
type OpenMeteoResponse struct {
	Daily struct {
		Time             []string  `json:"time"`
		MaxTemperature   []float64 `json:"temperature_2m_max"`
		MinTemperature   []float64 `json:"temperature_2m_min"`
		PrecipitationSum []float64 `json:"precipitation_sum"`
		WeatherCode      []int     `json:"weathercode"`
	} `json:"daily"`
}

// Fetch weather data from Open-Meteo API
func fetchWeatherData(latitude, longitude float64) ([]WeatherDay, error) {
	// Get current date
	now := time.Now()

	// Format API URL
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,weathercode&timezone=auto",
		latitude, longitude,
	)

	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Decode response
	var apiResp OpenMeteoResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	// Build weather days from API response
	var days []WeatherDay
	for i, dateStr := range apiResp.Daily.Time {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// Convert data
		weatherDay := WeatherDay{
			Date:          dateStr,
			FormattedDate: date.Format("Monday, Jan 2"),
			MaxTemp:       apiResp.Daily.MaxTemperature[i],
			MinTemp:       apiResp.Daily.MinTemperature[i],
			Precipitation: apiResp.Daily.PrecipitationSum[i],
			WeatherCode:   apiResp.Daily.WeatherCode[i],
			Icon:          getWeatherIcon(apiResp.Daily.WeatherCode[i]),
			IsWeekend:     date.Weekday() == time.Saturday || date.Weekday() == time.Sunday,
			IsToday:       date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day(),
		}

		days = append(days, weatherDay)
	}

	return days, nil
}

// Get weather icon based on weather code
func getWeatherIcon(code int) string {
	switch {
	case code == 0:
		return "sun" // Clear sky
	case code >= 1 && code <= 3:
		return "cloud-sun" // Mainly clear, partly cloudy, and overcast
	case code >= 45 && code <= 48:
		return "cloud-fog" // Fog and depositing rime fog
	case code >= 51 && code <= 55:
		return "cloud-drizzle" // Drizzle
	case code >= 61 && code <= 65:
		return "cloud-rain" // Rain
	case code >= 66 && code <= 67 || code >= 80 && code <= 82:
		return "cloud-snow" // Freezing rain or snow
	case code >= 71 && code <= 77 || code >= 85 && code <= 86:
		return "snowflake" // Snow fall
	case code >= 95 && code <= 99:
		return "cloud-lightning" // Thunderstorm
	default:
		return "question-circle"
	}
}

// Popular hiking parks with their coordinates
var PopularParks = map[string]Park{
	"Adirondack Park": {
		Name:        "Adirondack Park",
		Location:    "New York",
		Latitude:    44.125,
		Longitude:   -73.925,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e2/Adirondacks_View_from_Coney_Mtn.jpg/1024px-Adirondacks_View_from_Coney_Mtn.jpg",
		Description: "The largest protected area in the contiguous United States, with over 6 million acres of forests, lakes, and mountains.",
		License:     "CC BY-SA 4.0",
		Attribution: "Mwanner via Wikimedia Commons",
	},
	"Yosemite National Park": {
		Name:        "Yosemite National Park",
		Location:    "California",
		Latitude:    37.8651,
		Longitude:   -119.5383,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d6/Half_Dome_from_Glacier_Point%2C_Yosemite_NP_-_Diliff.jpg/1024px-Half_Dome_from_Glacier_Point%2C_Yosemite_NP_-_Diliff.jpg",
		Description: "Known for its waterfalls, deep valleys, grand meadows, and ancient giant sequoias.",
		License:     "CC BY-SA 3.0",
		Attribution: "Diliff via Wikimedia Commons",
	},
	"Grand Canyon National Park": {
		Name:        "Grand Canyon National Park",
		Location:    "Arizona",
		Latitude:    36.1069,
		Longitude:   -112.1129,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Grand_Canyon_view_from_Pima_Point_2010.jpg/1024px-Grand_Canyon_view_from_Pima_Point_2010.jpg",
		Description: "The Grand Canyon is 277 miles long, up to 18 miles wide and reaches a depth of over a mile.",
		License:     "CC BY-SA 3.0",
		Attribution: "Chensiyuan via Wikimedia Commons",
	},
	"Great Smoky Mountains National Park": {
		Name:        "Great Smoky Mountains National Park",
		Location:    "Tennessee",
		Latitude:    35.6131,
		Longitude:   -83.5532,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e8/Smoky_Mountains_-_scenic_landscape_-_near_Newfound_Gap.JPG/1024px-Smoky_Mountains_-_scenic_landscape_-_near_Newfound_Gap.JPG",
		Description: "America's most visited national park, known for its diverse plant and animal life and the beauty of its ancient mountains.",
		License:     "CC BY-SA 3.0",
		Attribution: "Brian Stansberry via Wikimedia Commons",
	},
	"Zion National Park": {
		Name:        "Zion National Park",
		Location:    "Utah",
		Latitude:    37.2982,
		Longitude:   -113.0263,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/Angels_Landing.jpg/1024px-Angels_Landing.jpg",
		Description: "Known for its steep red cliffs, emerald pools, and narrow slot canyons.",
		License:     "CC BY-SA 4.0",
		Attribution: "Diliff via Wikimedia Commons",
	},
	"Acadia National Park": {
		Name:        "Acadia National Park",
		Location:    "Maine",
		Latitude:    44.3386,
		Longitude:   -68.2733,
		ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Bass_Harbor_Head_Light_Station_2016.jpg/1024px-Bass_Harbor_Head_Light_Station_2016.jpg",
		Description: "The first national park east of the Mississippi River, featuring the highest rocky headlands along the Atlantic coastline.",
		License:     "Public Domain",
		Attribution: "National Park Service",
	},
}

// Sample trail data for each park, using open images
var parkTrails = map[string][]Trail{
	"Adirondack Park": {
		{
			Name:        "High Peaks Wilderness Loop",
			Length:      19.3,
			Difficulty:  "Difficult",
			Description: "A challenging loop through the High Peaks region offering spectacular views.",
			Elevation:   1463,
			EstTime:     "8-10 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/Adirondack_Loj.jpg/1024px-Adirondack_Loj.jpg",
			License:     "CC BY-SA 3.0",
			Attribution: "Mwanner via Wikimedia Commons",
		},
		{
			Name:        "Cascade Mountain Trail",
			Length:      7.4,
			Difficulty:  "Moderate",
			Description: "One of the most popular trails in the Adirondacks with panoramic summit views.",
			Elevation:   863,
			EstTime:     "4-5 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/9/94/Cascade_Mountain_NY_2.jpg/1024px-Cascade_Mountain_NY_2.jpg",
			License:     "CC BY-SA 4.0",
			Attribution: "Carl Heilman II via Wikimedia Commons",
		},
	},
	"Yosemite National Park": {
		{
			Name:        "Half Dome Trail",
			Length:      23.0,
			Difficulty:  "Very Difficult",
			Description: "Iconic trail with cable section to reach the summit of Half Dome.",
			Elevation:   1573,
			EstTime:     "10-12 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/9/96/Half_Dome_Cables.JPG/1024px-Half_Dome_Cables.JPG",
			License:     "CC BY-SA 3.0",
			Attribution: "Inklein via Wikimedia Commons",
		},
		{
			Name:        "Mist Trail to Vernal Fall",
			Length:      5.4,
			Difficulty:  "Moderate",
			Description: "Popular trail along the Merced River to the spectacular Vernal Fall.",
			Elevation:   305,
			EstTime:     "3 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/Vernal_fall_rainbow.jpg/1024px-Vernal_fall_rainbow.jpg",
			License:     "CC BY-SA 4.0",
			Attribution: "King of Hearts via Wikimedia Commons",
		},
	},
	"Grand Canyon National Park": {
		{
			Name:        "Bright Angel Trail",
			Length:      15.6,
			Difficulty:  "Difficult",
			Description: "The park's most popular trail into the canyon with regular rest houses.",
			Elevation:   1335,
			EstTime:     "8-10 hours round trip to Plateau Point",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/2/2a/Grand_Canyon_-_panoramio_%2830%29.jpg/1024px-Grand_Canyon_-_panoramio_%2830%29.jpg",
			License:     "CC BY 3.0",
			Attribution: "Gottfried Achenwall via Wikimedia Commons",
		},
		{
			Name:        "South Kaibab Trail",
			Length:      11.3,
			Difficulty:  "Difficult",
			Description: "Steep trail with expansive views and no water available.",
			Elevation:   1441,
			EstTime:     "6-8 hours round trip to Skeleton Point",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/South_Kaibab_Trail%2C_Grand_Canyon.jpg/1024px-South_Kaibab_Trail%2C_Grand_Canyon.jpg",
			License:     "CC BY 2.0",
			Attribution: "Michael Quinn via Wikimedia Commons",
		},
	},
	"Great Smoky Mountains National Park": {
		{
			Name:        "Alum Cave Trail to Mount LeConte",
			Length:      17.0,
			Difficulty:  "Moderate to Difficult",
			Description: "Scenic trail passing through Alum Cave to reach Mount LeConte summit.",
			Elevation:   763,
			EstTime:     "6-8 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8a/Alum_Cave_Bluffs_in_the_Great_Smoky_Mountains.jpg/1024px-Alum_Cave_Bluffs_in_the_Great_Smoky_Mountains.jpg",
			License:     "CC BY-SA 3.0",
			Attribution: "Brian Stansberry via Wikimedia Commons",
		},
		{
			Name:        "Chimney Tops Trail",
			Length:      6.1,
			Difficulty:  "Moderate to Difficult",
			Description: "Short but steep trail leading to rocky pinnacle with 360-degree views.",
			Elevation:   467,
			EstTime:     "3-4 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Chimney_Tops_Trail_%2850550822493%29.jpg/1024px-Chimney_Tops_Trail_%2850550822493%29.jpg",
			License:     "CC BY 2.0",
			Attribution: "Ken Lund via Wikimedia Commons",
		},
	},
	"Zion National Park": {
		{
			Name:        "Angels Landing",
			Length:      8.7,
			Difficulty:  "Difficult",
			Description: "Famous trail with chain sections along narrow ridges to a stunning viewpoint.",
			Elevation:   453,
			EstTime:     "4-5 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ae/Angels_Landing_Chain_Section.jpg/1024px-Angels_Landing_Chain_Section.jpg",
			License:     "CC BY-SA 4.0",
			Attribution: "Doc Searls via Wikimedia Commons",
		},
		{
			Name:        "The Narrows",
			Length:      15.0,
			Difficulty:  "Moderate to Difficult",
			Description: "Unique hiking experience through the narrowest section of Zion Canyon in the Virgin River.",
			Elevation:   334,
			EstTime:     "6-8 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3b/The_Narrows_in_Zion_Canyon.jpg/1024px-The_Narrows_in_Zion_Canyon.jpg",
			License:     "CC BY-SA 3.0",
			Attribution: "Jon Jasper via Wikimedia Commons",
		},
	},
	"Acadia National Park": {
		{
			Name:        "Precipice Trail",
			Length:      3.2,
			Difficulty:  "Very Difficult",
			Description: "Iron-rung route up the vertical eastern face of Champlain Mountain.",
			Elevation:   305,
			EstTime:     "2-3 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/5/5e/Precipice_Trail_%28Acadia_National_Park%29.jpg/1024px-Precipice_Trail_%28Acadia_National_Park%29.jpg",
			License:     "Public Domain",
			Attribution: "U.S. National Park Service",
		},
		{
			Name:        "Jordan Pond Path",
			Length:      5.1,
			Difficulty:  "Easy",
			Description: "Scenic loop around Jordan Pond with views of the Bubbles mountains.",
			Elevation:   30,
			EstTime:     "2-3 hours",
			ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0b/Jordan_Pond.JPG/1024px-Jordan_Pond.JPG",
			License:     "Public Domain",
			Attribution: "U.S. National Park Service",
		},
	},
}

// Get trails for a specific park
func getTrailsForPark(parkName string, lat, lng float64) []Trail {
	// Initialize random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// If we have predefined trails for this park, use them
	if trails, ok := parkTrails[parkName]; ok {
		// Add slight randomization to coordinates and recommended day
		for i := range trails {
			// Add small offsets to start positions
			trails[i].StartLat = lat + (r.Float64()*0.03 - 0.015)
			trails[i].StartLng = lng + (r.Float64()*0.03 - 0.015)

			// Recommend a day (usually within the next 7 days)
			trails[i].RecommendDay = r.Intn(7)
		}
		return trails
	}

	// Otherwise, generate some generic trails with nature-related names
	var trailNames = []string{
		"Ridge Loop", "Valley View", "Summit Trail", "Waterfall Path",
		"Mountain Circuit", "Forest Trail", "Lake View Trail", "River Walk",
		"Wildflower Trail", "Eagle Point", "Canyon Trek", "Meadow Loop",
	}

	var difficulties = []string{"Easy", "Moderate", "Difficult", "Very Difficult"}

	var descriptions = []string{
		"Beautiful trail with views of surrounding landscapes.",
		"Popular path through diverse terrain with plenty of wildlife.",
		"Challenging hike with rewarding panoramic vistas at the summit.",
		"Scenic route following a stream to a secluded viewpoint.",
		"Family-friendly trail through meadows and forests.",
		"Diverse trail featuring waterfalls and mountain views.",
	}

	// Open-source nature images from Wikimedia Commons
	var images = []string{
		"https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Alpine_trail_in_Switzerland.jpg/1024px-Alpine_trail_in_Switzerland.jpg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/a/ab/Forest_trail_in_Thetford_Forest_-_geograph.org.uk_-_1403862.jpg/1024px-Forest_trail_in_Thetford_Forest_-_geograph.org.uk_-_1403862.jpg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/5/5c/Trail_in_the_Forest_%28250773707%29.jpeg/1024px-Trail_in_the_Forest_%28250773707%29.jpeg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/3/3a/Autumn_forest_path.jpg/1024px-Autumn_forest_path.jpg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/c/c5/Joshua_Tree_Hiking_Trail.jpg/1024px-Joshua_Tree_Hiking_Trail.jpg",
	}

	var licenses = []string{
		"CC BY-SA 4.0", "CC BY 3.0", "CC BY-SA 3.0", "Public Domain",
	}

	var attributions = []string{
		"Mark Janes via Wikimedia Commons",
		"Forest Service, USDA",
		"Wikimedia Commons Contributors",
		"Public Domain",
		"Ashley Whitworth via Wikimedia Commons",
	}

	// Generate 2-4 trails
	numTrails := r.Intn(3) + 2
	var trails []Trail

	for i := 0; i < numTrails; i++ {
		// Generate random trail properties
		name := trailNames[r.Intn(len(trailNames))]
		if len(parkName) > 0 {
			name = parkName + " " + name
		}

		// Avoid duplicate names
		for j := range trails {
			if trails[j].Name == name {
				name = name + " " + string('A'+i)
				break
			}
		}

		length := 1.5 + r.Float64()*20.0
		diff := difficulties[r.Intn(len(difficulties))]
		desc := descriptions[r.Intn(len(descriptions))]
		elevation := 50.0 + r.Float64()*1500.0

		// Estimate time based on length and difficulty
		var estTime string
		hours := int(length / 3.0)
		if hours < 1 {
			estTime = "30-45 minutes"
		} else if hours < 2 {
			estTime = "1-2 hours"
		} else {
			estTime = fmt.Sprintf("%d-%d hours", hours, hours+2)
		}

		imgURL := images[r.Intn(len(images))]
		license := licenses[r.Intn(len(licenses))]
		attribution := attributions[r.Intn(len(attributions))]

		// Create trail with small position offset from park center
		trails = append(trails, Trail{
			Name:         name,
			Length:       length,
			Difficulty:   diff,
			Description:  desc,
			Elevation:    elevation,
			EstTime:      estTime,
			ImageURL:     imgURL,
			License:      license,
			Attribution:  attribution,
			StartLat:     lat + (r.Float64()*0.03 - 0.015),
			StartLng:     lng + (r.Float64()*0.03 - 0.015),
			RecommendDay: r.Intn(7),
		})
	}

	return trails
}
