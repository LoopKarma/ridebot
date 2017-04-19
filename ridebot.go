package main

import (
	"fmt"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"time"
)

var targetPoint = Coordinates{
	Lat:52.4833,
	Lon:13.3833,
}
var apiHost = "http://api.openweathermap.org/data/2.5/"
var apiKey = "b2896361aa576c90279d9f966b80aa8d"


func main() {
	var ride = RideParams{
		duration:4,
		isNightRide:true,
		minTemp:-5,
		isRainRide:true,
	}

	forecast := getImportantForecast(ride)
	if (len(forecast.List) > 0) {
		message := "\n\n" + forecast.City.Name + ": \n";
		for _,point := range forecast.List {
			message += "In " + point.DtTxt + " min = " + floatToString(point.Main.TempMin, 4) + "˚; max = " + floatToString(point.Main.TempMax) + "˚\n"
			message += point.Weather[0].Description + "\n"
		}
		fmt.Println("\n\n\n")
		fmt.Println(message)
		fmt.Println("\n\n\n")
	}
}

func getImportantForecast(ride RideParams) ForecastResponse {
	weather := getCurrentWeather()
	forecast := getForecast()



	var importantPoints []List

	for _, point := range forecast.List {
		if isImportantPoint(weather, point, ride) {
			importantPoints = append(importantPoints, point)
		}
	}
	forecast.List = importantPoints


	forecast.Sys = weather.Sys

	return forecast
}

func isImportantPoint(currentWeather WeatherResponse, point List, ride RideParams) bool {
	sunset := currentWeather.Sys.Sunset
	rideEnd := currentWeather.Dt + (ride.duration * 3600)

	//fmt.Println("\n sunset ", timestampToHumanType(sunset))
	//fmt.Println("human time of point ", point.DtTxt)

	if ride.isNightRide == false && point.Dt > sunset { //point in the evening
		//fmt.Println("point.Dt > sunset ", point.Dt > sunset)
		return false
	}


	//fmt.Println("\n curtime ", timestampToHumanType(currentWeather.Dt))
	//fmt.Println("currentWeather.Dt > point.Dt ", currentWeather.Dt > point.Dt)

	if currentWeather.Dt > point.Dt {     //point in the past
		return isGoodConditions(ride, point)
	}

	fmt.Println("\n point.Dt ",point.DtTxt)
	fmt.Println("\n point.Dt ",point.Dt)
	fmt.Println("rideEnd ",timestampToHumanType(rideEnd))
	fmt.Println("rideEnd ",rideEnd)
	fmt.Println("point.Dt < rideEnd", point.Dt < rideEnd)

	if point.Dt < rideEnd {
		return isGoodConditions(ride, point)
	}


	return false
}

func isGoodConditions(ride RideParams, point List) bool {
	goodWeatherConditions := map[int]string {
		800:"clear sky",
		801:"few clouds",
		802:"scattered clouds",
		803:"broken clouds",
		804:"overcast clouds",
	}
	if ride.isRainRide {
		goodWeatherConditions[500] = "light rain"
	}
	for key, _ := range goodWeatherConditions {
		if point.Weather[0].ID == key {
			return ride.minTemp <= int(point.Main.TempMin)
		}
	}

	return false
}


func getForecast() ForecastResponse {
	requestParams := map[string]string{
		"lat": floatToString(targetPoint.Lat),
		"lon":   floatToString(targetPoint.Lon),
	}

	resp, err := http.Get(createApiUri("forecast", requestParams))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := ForecastResponse{}


	if err := json.Unmarshal(content, &result); err != nil {
		panic(err)
	}

	//TODO consider to remove debug messages
	//stringContent := string(content)
	//fmt.Println("\n\n\n")
	//fmt.Println(stringContent)
	//fmt.Println("\n\n\n")

	return result
}

func getCurrentWeather() WeatherResponse {
	requestParams := map[string]string{
		"lat": floatToString(targetPoint.Lat),
		"lon":   floatToString(targetPoint.Lon),
	}

	resp, err := http.Get(createApiUri("weather", requestParams))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := WeatherResponse{}


	if err := json.Unmarshal(content, &result); err != nil {
		panic(err)
	}

	//TODO consider to remove debug messages
	//stringContent := string(content)
	//fmt.Println("\n\n\n")
	//fmt.Println(stringContent)
	//fmt.Println("\n\n\n")

	return result
}

func createApiUri(method string, parameters map[string]string) string {
	values := url.Values{}
	for key, value := range parameters {
		values.Add(key, value)
	}
	values.Add("APPID", apiKey)
	values.Add("units", "metric")
	values.Add("lang", "ru")

	uri := apiHost + method + "?" + values.Encode()
	fmt.Println(uri)

	return uri
}

func floatToString(input_num float64, precision ...int) string {
	var prec int
	if precision == nil {
		prec = 2
	} else {
		prec = precision[0]
	}
	return strconv.FormatFloat(input_num, 'f', prec, 64)
}

func timestampToHumanType(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}


type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ForecastResponse struct {
	Cod string `json:"cod"`
	Message float64 `json:"message"`
	Cnt int `json:"cnt"`
	List []List `json:"list"`
	City struct {
		    ID int `json:"id"`
		    Name string `json:"name"`
		    Coord Coordinates `json:"coord"`
		    Country string `json:"country"`
	    } `json:"city"`
	Sys Sys
}

type WeatherResponse struct {
	Base   string `json:"base"`
	Clouds struct {
		       All int `json:"all"`
	       } `json:"clouds"`
	Cod   int `json:"cod"`
	Coord Coordinates `json:"coord"`
	Dt   int `json:"dt"`
	ID   int `json:"id"`
	Main struct {
		       Humidity int     `json:"humidity"`
		       Pressure int     `json:"pressure"`
		       Temp     float64 `json:"temp"`
		       TempMax  float64 `json:"temp_max"`
		       TempMin  float64 `json:"temp_min"`
	       } `json:"main"`
	Name string `json:"name"`
	Sys Sys `json:"sys"`
	Visibility int `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ID          int    `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		       Deg   int     `json:"deg"`
		       Speed float64 `json:"speed"`
	       } `json:"wind"`
}

type List struct {
	Dt int `json:"dt"`
	Main struct {
		   Temp float64 `json:"temp"`
		   TempMin float64 `json:"temp_min"`
		   TempMax float64 `json:"temp_max"`
		   Pressure float64 `json:"pressure"`
		   SeaLevel float64 `json:"sea_level"`
		   GrndLevel float64 `json:"grnd_level"`
		   Humidity int `json:"humidity"`
		   TempKf float64 `json:"temp_kf"`
	   } `json:"main"`
	Weather []struct {
		ID int `json:"id"`
		Main string `json:"main"`
		Description string `json:"description"`
		Icon string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		   All int `json:"all"`
	   } `json:"clouds"`
	Wind struct {
		   Speed float64 `json:"speed"`
		   Deg float64 `json:"deg"`
	   } `json:"wind"`
	Rain struct {
		   ThreeH float64 `json:"3h"`
	   } `json:"rain"`
	Sys struct {
		   Pod string `json:"pod"`
	   } `json:"sys"`
	DtTxt string `json:"dt_txt"`
	Snow struct {
		   ThreeH float64 `json:"3h"`
	   } `json:"snow,omitempty"`
}

type Sys struct {
	Country string  `json:"country"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
	Type    int     `json:"type"`
}

type RideParams struct {
	duration    int
	isNightRide bool
	minTemp     int
	isRainRide  bool
}
