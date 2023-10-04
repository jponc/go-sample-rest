package openmateo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-sample-rest/internal/types"

	log "github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	httpClient HTTPClient
}

func NewClient(httpClient HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetLatestWeatherData(latitude float64, longitude float64) (*types.WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", latitude, longitude)

	log.Infof(fmt.Sprintf("Requesting: %s", url))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request weather data: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest weather, response status code: %d", resp.StatusCode)
	}

	var body OpenMateoForecastResponseBody

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &types.WeatherData{
		Latitude:      body.Latitude,
		Longitude:     body.Longitude,
		Temperature:   body.CurrentWeather.Temperature,
		WindDirection: body.CurrentWeather.WindDirection,
		WindSpeed:     body.CurrentWeather.WindSpeed,
	}, nil
}
