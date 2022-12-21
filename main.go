package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var lat string
var long string

func getWeather(url string) []string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "AOPRIVTEMP/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type Weather struct {
		Properties struct {
			Timeseries []struct {
				Time time.Time `json:"time"`
				Data struct {
					Instant struct {
						Details struct {
							AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
							AirTemperature        float64 `json:"air_temperature"`
							CloudAreaFraction     float64 `json:"cloud_area_fraction"`
							RelativeHumidity      float64 `json:"relative_humidity"`
							WindFromDirection     float64 `json:"wind_from_direction"`
							WindSpeed             float64 `json:"wind_speed"`
						} `json:"details"`
					} `json:"instant"`
				} `json:"data,omitempty"`
			} `json:"timeseries"`
		} `json:"properties"`
	}

	var weather Weather

	err = json.Unmarshal(bytes, &weather)
	if err != nil {
		fmt.Println("error:", err)
	}

	prefix := "yr"

	var metricsOut []string

	// metricsOut = append(metricsOut, fmt.Sprintf("%v_%v %v", prefix, "time", weather.Properties.Timeseries[1].Time))
	metricsOut = append(metricsOut, fmt.Sprintf("%v_%v %.1f", prefix, "air_temperature", weather.Properties.Timeseries[1].Data.Instant.Details.AirTemperature))
	metricsOut = append(metricsOut, fmt.Sprintf("%v_%v %.1f", prefix, "relative_humidity", weather.Properties.Timeseries[1].Data.Instant.Details.RelativeHumidity))

	return metricsOut

}

func MetricsHttp(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s", lat, long)

	allOut := getWeather(url)
	fmt.Fprintln(w, strings.Join(allOut, "\n"))
}

func main() {
	if os.Getenv("YR_LAT") == "" || os.Getenv("YR_LONG") == "" {
		fmt.Println("MISSING ENV VARS...")
		os.Exit(0)
	} else {
		lat = os.Getenv("YR_LAT")
		long = os.Getenv("YR_LONG")
	}
	port := "9118"
	http.HandleFunc("/metrics", MetricsHttp)
	panic(http.ListenAndServe("0.0.0.0:"+port, nil))
}
