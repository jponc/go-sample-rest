package types

type WeatherData struct {
	Id            string  `json:"id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Temperature   float64 `json:"temperature"`
	WindDirection float64 `json:"wind_direction"`
	WindSpeed     float64 `json:"wind_speed"`
	CreatedAt     string  `json:"created_at"`
}

type GetLatestWeatherResponse WeatherData

type GetWeatherHistoryResponse []WeatherData
