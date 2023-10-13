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
				TimeEpoch    int     `json:"time_epoch"`
				TempC        float64 `json:"temp_c"`
				ChanceOfRain float64 `json:"chance_of_rain"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	country := "London"
	if len(os.Args) >= 2 {
		country = os.Args[1]
	}
	var endpoint = "http://api.weatherapi.com/v1/forecast.json?key=4bec63a552834673a8d101342231310&q=" + country + "&days=1&aqi=no&alerts=no"
	res, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("Open Weather Api is not responding")
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
	Current, location, hours := weather.Current, weather.Location, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%v,%v %.0fC and weather is %v\n",
		location.Name, location.Country, Current.TempC, Current.Condition.Text,
	)
	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)
		if date.Before(time.Now()) {
			continue
		}
		message := fmt.Sprintf("%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
		if hour.ChanceOfRain < 40 {
			fmt.Println(message)
		} else {
			color.Red(message)
		}
	}
}
