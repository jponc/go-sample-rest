package openmateo

type OpenMateoCurrentWeather struct {
	Temperature   float64 `json:"temperature"`
	WindDirection float64 `json:"winddirection"`
	WindSpeed     float64 `json:"windspeed"`
}

type OpenMateoForecastResponseBody struct {
	Latitude       float64                 `json:"latitude"`
	Longitude      float64                 `json:"longitude"`
	CurrentWeather OpenMateoCurrentWeather `json:"current_weather"`
}
