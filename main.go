package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Astana"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=000eed2ab597491185432211231010&q=" + q + "&days=1&aqi=yes&alerts=yes")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)
	for _, h := range hours {
		date := time.Unix(h.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		message := fmt.Sprintf("%s - %.0fC, %.0f%%, %s\n", date.Format("15:04"), h.TempC, h.ChanceOfRain, h.Condition.Text)
		if h.ChanceOfRain < 40 {
			color.Green(message)
		} else {
			color.Red(message)
		}
	}
}
