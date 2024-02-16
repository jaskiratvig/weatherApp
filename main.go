package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Struct for simplified weather response
type WeatherResponse struct {
	Condition   string `json:"condition"`
	Temperature string `json:"temperature"`
}

// Struct to decode error response from OpenWeatherMap
type APIErrorResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

// Handler for the weather endpoint
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key for OpenWeatherMap is not set", http.StatusInternalServerError)
		return
	}

	log.Println("Using OpenWeather API Key:", apiKey)

	// Extract query parameters for latitude and longitude
	query := r.URL.Query()
	lat := query.Get("lat")
	lon := query.Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "Please provide both latitude and longitude", http.StatusBadRequest)
		return
	}

	// Construct the API request URL
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to request weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check for invalid API key
	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("Invalid API key detected")
		var apiError APIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			http.Error(w, fmt.Sprintf("Invalid API key: %s", apiError.Message), http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
		}
		return
	}

	// Read and log the entire API response body for debugging
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read the weather API response body: %v", err)
		http.Error(w, "Failed to process weather data", http.StatusInternalServerError)
		return
	}

	// Log the response body for inspection
	log.Printf("API Response: %s", string(bodyBytes))

	// Convert the response body back to an io.Reader for JSON decoding
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Attempt to parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		http.Error(w, "Failed to decode weather data", http.StatusInternalServerError)
		return
	}

	// Simplify the response based on weather condition and temperature

	// Safely extract weather condition and temperature
	weather, weatherExists := result["weather"].([]interface{})
	var weatherCondition string
	if weatherExists && len(weather) > 0 {
		condition, conditionExists := weather[0].(map[string]interface{})["main"].(string)
		if conditionExists {
			weatherCondition = condition
		} else {
			weatherCondition = "Unknown"
		}
	} else {
		weatherCondition = "Unknown"
	}

	// Safely extract temperature
	main, mainExists := result["main"].(map[string]interface{})
	var temp float64
	if mainExists {
		tempValue, tempExists := main["temp"].(float64)
		if tempExists {
			temp = tempValue
		} else {
			// Handle the case where the temperature is not available/parseable
			http.Error(w, "Temperature data is unavailable", http.StatusInternalServerError)
			return
		}
	} else {
		// Handle the case where the "main" key is missing or is not the expected type
		http.Error(w, "Weather data format is unexpected", http.StatusInternalServerError)
		return
	}

	// Determine temperature description
	var tempDesc string
	switch {
	case temp < 10:
		tempDesc = "Cold"
	case temp >= 10 && temp <= 25:
		tempDesc = "Moderate"
	case temp > 25:
		tempDesc = "Hot"
	}

	// Respond to the client
	response := WeatherResponse{
		Condition:   strings.ToLower(weatherCondition),
		Temperature: tempDesc,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/weather", weatherHandler)

	port := ":8080"
	fmt.Println("Server listening on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
