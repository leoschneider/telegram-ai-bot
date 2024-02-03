package weather_api

const WeatherUiEndpoint = "https://www.weatherapi.com/weather/?q="

var emojiCodeMap = map[int]string{
	1000: "☀️", // "Sunny"
	1003: "🌤",  // "Partly Cloudy"
	1006: "⛅",  // "Cloudy"
	1009: "☁️", // "Overcast"
	1030: "🌫",  // "Mist"
	1063: "🌧️", // "Patchy rain nearby"
	1066: "🌨",  // "Patchy snow nearby"
	1069: "🌤️", // "Patchy sleet nearby"
	1072: "❄️", // "Patchy freezing drizzle nearby"
	1087: "⚡ ", // "Thundery outbreaks in nearby"
	1114: "🌨",  // "Blowing snow"
	1117: "❄️", // "Blizzard"
	1135: "🌫",  // "Fog"
	1147: "❄️", // "Freezing fog"
	1150: "🌧️", // "Patchy light drizzle"
	1153: "🌧️", // "Light drizzle"
	1168: "❄️", // "Freezing drizzle"
	1171: "❄️", // "Heavy freezing drizzle"
	1180: "🌧️", // "Patchy light rain"
	1183: "🌧️", // "Light rain"
	1186: "🌧️", // "Moderate rain at times"
	1189: "🌧️", // "Moderate rain"
	1192: "🌧️", // "Heavy rain at times"
	1195: "🌧️", // "Heavy rain"
	1198: "❄️", // "Light freezing rain"
	1201: "❄️", // "Moderate or heavy freezing rain"
	1204: "🌤",  // "Light sleet"
	1207: "🌤",  // "Moderate or heavy sleet"
	1210: "❄️", // "Patchy light snow"
	1213: "❄️", // "Light snow"
	1216: "❄️", // "Patchy moderate snow"
	1219: "❄️", // "Moderate snow"
	1222: "❄️", // "Patchy heavy snow"
	1225: "❄️", // "Heavy snow"
	1237: "❄️", // "Ice pellets"
	1240: "🌧️", // "Light rain shower"
	1243: "☔ ", // "Moderate or heavy rain shower"
	1246: "☔ ", // "Torrential rain shower"
	1249: "🌧️", // "Light sleet showers"
	1252: "🌧️", // "Moderate or heavy sleet showers"
	1255: "🌨",  // "Light snow showers"
	1258: "🌨",  // "Moderate or heavy snow showers"
	1261: "❄️", // "Light showers of ice pellets"
	1264: "❄️", // "Moderate or heavy showers of ice pellets"
	1273: "⛈️", // "Patchy light rain in area with thunder"
	1276: "⛈️", // "Moderate or heavy rain in area with thunder"
	1279: "⛈️", // "Patchy light snow in area with thunder"
	1282: "⛈️", // "Moderate or heavy snow in area with thunder"
}

type Forecastday struct {
	Hour []WeatherStatus `json:"hour"`
}

type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

type Weather struct {
	Location Location      `json:"location"`
	Current  WeatherStatus `json:"current"`
	Forecast Forecast      `json:"forecast"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TimezoneID     string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type WeatherStatus struct {
	LastUpdatedEpoch int64      `json:"last_updated_epoch"`
	LastUpdated      string     `json:"last_updated"`
	TempC            float64    `json:"temp_c"`
	TempF            float64    `json:"temp_f"`
	IsDay            int        `json:"is_day"`
	Condition        Condition  `json:"condition"`
	WindMph          float64    `json:"wind_mph"`
	WindKph          float64    `json:"wind_kph"`
	WindDegree       int        `json:"wind_degree"`
	WindDir          string     `json:"wind_dir"`
	PressureMb       int        `json:"pressure_mb"`
	PressureIn       float64    `json:"pressure_in"`
	PrecipMm         float64    `json:"precip_mm"`
	PrecipIn         float64    `json:"precip_in"`
	Humidity         int        `json:"humidity"`
	Cloud            int        `json:"cloud"`
	FeelsLikeC       float64    `json:"feelslike_c"`
	FeelsLikeF       float64    `json:"feelslike_f"`
	VisKm            int        `json:"vis_km"`
	VisMiles         int        `json:"vis_miles"`
	UV               int        `json:"uv"`
	GustMph          float64    `json:"gust_mph"`
	GustKph          float64    `json:"gust_kph"`
	AirQuality       AirQuality `json:"air_quality"`
	Time             string     `json:"time"`
	TimeEpoch        int64      `json:"time_epoch"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type AirQuality struct {
	CO           float64 `json:"co"`
	NO2          float64 `json:"no2"`
	O3           float64 `json:"o3"`
	SO2          int     `json:"so2"`
	PM2_5        float64 `json:"pm2_5"`
	PM10         int     `json:"pm10"`
	USEpaIndex   int     `json:"us-epa-index"`
	GbDefraIndex int     `json:"gb-defra-index"`
}
