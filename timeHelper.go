package main

import (
	"time"
	"net/http"
	"net/url"
	"fmt"
)

var apiHost = "https://maps.googleapis.com/"
var apiMethod = "maps/api/timezone/json"

func getLocalTimeByCoordinates(coord Coordinates) time.Time {
	values := url.Values{}
	parameters := map[string]string {
		"location": string(coord.Lat)+","+string(coord.Lon),
		"timestamp": string(time.Now().Unix()),
	}
	for key, value := range parameters {
		values.Add(key, value)
	}
	requestUrl = apiHost + apiMethod + "?" + values.Encode()


}



