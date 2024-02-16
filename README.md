# Weather API Server

This project is a simple HTTP server written in Go that uses the OpenWeatherMap API to fetch current weather conditions based on latitude and longitude coordinates provided by the client. It returns a simplified response indicating the general weather condition (e.g., snow, rain) and a qualitative temperature assessment (hot, cold, or moderate).

## Prerequisites

Before you can run this server, you need to obtain an API key from OpenWeatherMap by following these steps:

1. Sign up for an account at [OpenWeatherMap](https://openweathermap.org/api).
2. Navigate to the API keys section and generate a new API key.
3. Once you have your API key, you will need to set it as an environment variable on your system.

## Setting Up Your Environment

To run this project, ensure you have Go installed on your machine. If not, download it from the [official Go website](https://golang.org/dl/).

## Running the Server

To get the server up and running, follow these steps:

1. Clone this repository to your local machine.
2. Open a terminal and navigate to the directory where you cloned the repository.
3. Set your OpenWeatherMap API key as an environment variable using the following command:
   export OPENWEATHER_API_KEY=your_api_key_here

Run the server with the command:

go run main.go

The server will start listening on port 8080. Access the weather endpoint by sending a GET request to http://localhost:8080/weather with lat (latitude) and lon (longitude) query parameters.

Example request:

curl "http://localhost:8080/weather?lat=35.6895&lon=139.6917"