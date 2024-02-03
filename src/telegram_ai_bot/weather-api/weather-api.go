package weather_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetWeatherResponse(weatherType string, key string, query string) (text string) {
	weatherResult := sendApiHttpRequest(weatherType, key, query)
	return weatherResult.buildMarkdownTable()
}

func sendApiHttpRequest(weatherType string, key string, query string) Weather {
	url := "https://api.weatherapi.com/v1/" + weatherType + ".json?key=" + key + "&q=" + query
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Request failed: %s", err)
		return Weather{}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s", err)
		return Weather{}
	}

	var deserialized Weather

	err = json.Unmarshal(body, &deserialized)

	return deserialized
}

func (w Weather) buildMarkdownTable() string {
	markdownTable := "`\n" +
		"|   |  Time  | Temp | Real Feel | rain | humidity | Wind |\n" +
		"|---|--------|------|-----------|------|----------|------|\n"

	w.Current.Time = w.Location.Localtime

	var laterForecastHours []WeatherStatus
	if w.Forecast.Forecastday != nil {
		for _, weatherAtHour := range w.Forecast.Forecastday[0].Hour {
			if weatherAtHour.TimeEpoch > w.Location.LocaltimeEpoch {
				laterForecastHours = append(laterForecastHours, weatherAtHour)
			}
		}
	}

	allWeatherStatuses := []WeatherStatus{w.Current}
	allWeatherStatuses = append(allWeatherStatuses, laterForecastHours...)

	for _, row := range allWeatherStatuses {
		markdownTable += "|" +
			fitToMaxSpace(emojiCodeMap[row.Condition.Code], len(" ")) + " | " +
			fitToMaxSpace(timestampToTime(row.Time), len(" Time ")) + " | " +
			fitToMaxSpace(floatToString(row.TempC), len("Temp")) + " | " +
			fitToMaxSpace(floatToString(row.FeelsLikeC), len("Real Feel")) + " | " +
			fitToMaxSpace(floatToString(row.PrecipMm), len("rain")) + " | " +
			fitToMaxSpace(intToString(row.Humidity), len("humidity")) + " | " +
			fitToMaxSpace(floatToString(row.WindKph), len("Wind")) + " |\n"
	}
	markdownTable += "`"

	return markdownTable
}

func fitToMaxSpace(columnValue string, maxChars int) string {
	for len(columnValue) < maxChars {
		if maxChars-len(columnValue) == 1 {
			return columnValue + " "
		}
		columnValue = " " + columnValue + " "
	}
	return columnValue
}

func timestampToTime(timestamp string) string {
	return strings.Split(timestamp, " ")[1]
}

func floatToString(float float64) string {
	return strconv.FormatFloat(float, 'f', -1, 64)
}

func intToString(integer int) string {
	return strconv.FormatInt(int64(integer), 10)
}
