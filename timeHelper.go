package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const googleURI = "https://maps.googleapis.com/maps/api/timezone/json?location=%f,%f&timestamp=%d&sensor=false"

func RetrieveGoogleTimezone(latitude float64, longitude float64) (googleTimezone *GoogleTimezone, err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		err = fmt.Errorf("%v", r)
	//	}
	//}()
	timestamp := time.Now().UTC().Unix()

	uri := fmt.Sprintf(googleURI, latitude, longitude, timestamp)

	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if err != nil {
		return googleTimezone, err
	}



	// Convert the response to a byte array
	rawDocument, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return googleTimezone, err
	}

	// Unmarshal the response to a GoogleTimezone object
	googleTimezone = new(GoogleTimezone)
	if err = json.Unmarshal(rawDocument, googleTimezone); err != nil {
		return googleTimezone, err
	}

	if googleTimezone.Status != "OK" {
		err = fmt.Errorf("Error : Google Status : %s", googleTimezone.Status)
		return googleTimezone, err
	}

	//if len(GoogleTimezone.TimezoneID) == 0 {
	//	err = fmt.Errorf("Error : No Timezone Id Provided")
	//	return googleTimezone, err
	//}
	googleTimezone.Timestamp = timestamp + int64(googleTimezone.RawOffset)

	return googleTimezone, err
}

type GoogleTimezone struct {
	DstOffset    float64 `json:"dstOffset"`
	RawOffset    float64 `json:"rawOffset"`
	Status       string  `json:"status"`
	TimezoneID   string  `json:"timeZoneId"`
	TimezoneName string  `json:"timeZoneName"`
	Timestamp    int64
}
