package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type WeatherAPIResponse struct {
	Forecasts []Forecast `json:"list"`
}

type GeocodingAPIResponse struct {
	Locations []Location
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Forecast struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	}
}

func main() {
	//Have OPEN_WEATHER_MAP_API_KEY be set to api key
	apiKey := os.Getenv("OPEN_WEATHER_MAP_API_KEY")

	args := os.Args

	if len(args) < 2 {
		panic("Failed to provide location.")
	}

	location := args[1]

	lat, lon := getLatLon(location, apiKey)

	res, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=imperial&appid=%s", lat, lon, apiKey))

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	var apiResponse WeatherAPIResponse
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Here is the weather forecast for %s.\n", location)
	for _, forecaset := range apiResponse.Forecasts {
		fmt.Printf("On %s, Temp: %f°F, Feels Like: %f°F.\n", getFormattedDate(forecaset.Dt), forecaset.Main.Temp, forecaset.Main.FeelsLike)
	}
}

func getFormattedDate(dt int64) string {
	tm := time.Unix(dt, 0)
	formattedTm := tm.Format(time.RFC1123)
	return formattedTm
}

func getLatLon(loc string, key string) (float64, float64) {
	res, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", loc, key))

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var locations []Location

	json.Unmarshal(body, &locations)

	return locations[0].Lat, locations[0].Lon
}
