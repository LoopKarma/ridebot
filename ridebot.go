package main

import (
	"fmt"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"time"
	"gopkg.in/telegram-bot-api.v4"
)

var targetPoint = Coordinates{
	Lat:52.498126,
	Lon:13.399683,
}
var apiHost = "http://api.openweathermap.org/data/2.5/"
var apiKey = "b2896361aa576c90279d9f966b80aa8d"


func main() {
	//bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	var ride = RideParams{
		duration:5,
		isNightRide:false,
		minTemp:5,
		isRainRide:true,
		location:targetPoint,
		//debug: true,
	}
	fmt.Printf("Your ride params:\n%+v\n", ride)
	fmt.Println("**************************************************")

	forecast := getForecastForRide(ride)
	if (len(forecast.List) > 0) {
		fmt.Println(formatRideForecast(forecast))
	} else {
		fmt.Println("No forecast for your ride")
	}
}

func formatRideForecast(forecast ForecastResponse) string {
	headerFormat := "\n%s:\n"
	message := fmt.Sprintf(headerFormat, forecast.City.Name )
	for _,point := range forecast.List {
		pointFormat := "%s [min = %.2f˚; max = %.2f˚]\nDescription: %s. Result: %t\n"
		message += fmt.Sprintf(
			pointFormat,
			point.DtTxt,
			point.Main.TempMin,
			point.Main.TempMax,
			point.Weather[0].Description,
			point.isGoodConditions,
		)
	}
	return message
}

func getForecastForRide(ride RideParams) ForecastResponse {
	weather := getCurrentWeather(ride.location.Lat, ride.location.Lon)
	forecast := getForecast(ride.location.Lat, ride.location.Lon)

	var importantPoints []List
	var previousPoint List
	var timeBetweenTwoPoints int
	timestamp := int(time.Now().UTC().Unix())
	rideEnd := timestamp + int((ride.duration * 3600))
	sunset := weather.Sys.Sunset

	for _, point := range forecast.List {
		var add = false

		if (ride.debug) {
			fmt.Println("\n\ntimestamp point.Dt = ",point.Dt)
			fmt.Println("timestamp rideEnd = ",rideEnd)
			fmt.Println("Date point.Dt ",point.DtTxt)
			fmt.Println("Date rideEnd ", timestampToTimeUTC(rideEnd))
			fmt.Println("point.Dt < rideEnd", point.Dt < rideEnd)
			fmt.Println("point.Dt - rideEnd", point.Dt - rideEnd)
		}

		if ride.isNightRide == false && point.Dt > sunset {//point in the evening
			continue
		}
		if timestamp > point.Dt {//point in the past
			add = true
		}
		if point.Dt < rideEnd {//point before ride ending
			add = true
		}
		if &previousPoint != nil {
			timeBetweenTwoPoints = point.Dt - previousPoint.Dt
		}
		if &timeBetweenTwoPoints != nil  && (point.Dt - rideEnd) < timeBetweenTwoPoints / 2 {
			add = true
		}
		if add {
			point.isGoodConditions = isGoodConditions(ride, point)
			importantPoints = append(importantPoints, point)
		}
		previousPoint = point
	}
	forecast.List = importantPoints

	forecast.Sys = weather.Sys

	return forecast
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


func getForecast(latitude float64, longitude float64) ForecastResponse {
	requestParams := map[string]string{
		"lat": floatToString(latitude),
		"lon":   floatToString(longitude),
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

	return result
}

func getCurrentWeather(latitude float64, longitude float64) WeatherResponse {
	requestParams := map[string]string{
		"lat": floatToString(latitude),
		"lon":   floatToString(longitude),
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
	//fmt.Println(uri)
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

func timestampToTimeUTC(timestamp int) string {
	return time.Unix(int64(timestamp), 0).UTC().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
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
	isGoodConditions bool
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
	duration    float64
	isNightRide bool
	minTemp     int
	isRainRide  bool
	location    Coordinates
	debug       bool
}
