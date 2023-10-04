package weatherservice

import (
	"fmt"
	"net/http"
	"strconv"

	"go-sample-rest/internal/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	log "github.com/sirupsen/logrus"
)

type WeatherDataClient interface {
	GetLatestWeatherData(lat, long float64) (*types.WeatherData, error)
}

type WeatherDataRepository interface {
	GetLatestWeatherData(lat, long float64) (*types.WeatherData, error)
	GetWeatherHistory(lat, long float64) ([]*types.WeatherData, error)
	SaveWeatherData(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error)
}

type Service struct {
	weatherDataClient     WeatherDataClient
	weatherDataRepository WeatherDataRepository
}

func NewService(weatherDataClient WeatherDataClient, weatherDataRepository WeatherDataRepository) *Service {
	return &Service{
		weatherDataClient:     weatherDataClient,
		weatherDataRepository: weatherDataRepository,
	}
}

func (s *Service) GetLatestWeather(w http.ResponseWriter, r *http.Request) {
	lat, long, err := s.getLatLong(r)
	if err != nil {
		log.Errorf("failed to get lat long from request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weatherData, err := s.weatherDataRepository.GetLatestWeatherData(lat, long)
	if err != nil {
		log.Errorf("failed to get weather data from repository: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if weatherData == nil {
		log.Infof("weather data not found: lat(%f), long(%f)", lat, long)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.JSON(w, r, weatherData)
}

func (s *Service) GetWeatherHistory(w http.ResponseWriter, r *http.Request) {
	lat, long, err := s.getLatLong(r)
	if err != nil {
		log.Errorf("failed to get lat long from request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weatherHistory, err := s.weatherDataRepository.GetWeatherHistory(lat, long)
	if err != nil {
		log.Errorf("failed to get weather data history from repository: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, weatherHistory)
}

func (s *Service) UpdateWeather(w http.ResponseWriter, r *http.Request) {
	lat, long, err := s.getLatLong(r)
	if err != nil {
		log.Errorf("failed to get lat long from request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weatherData, err := s.weatherDataClient.GetLatestWeatherData(lat, long)
	if err != nil {
		log.Errorf("failed to get weather data from weather data client: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if weatherData == nil {
		log.Infof("weather data not found: lat(%f), long(%f)", lat, long)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	savedWeatherData, err := s.weatherDataRepository.SaveWeatherData(
		lat,
		long,
		weatherData.Temperature,
		weatherData.WindDirection,
		weatherData.WindSpeed,
	)
	if err != nil {
		log.Errorf("failed to save weather data to repository: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, savedWeatherData)
}

func (s *Service) getLatLong(r *http.Request) (float64, float64, error) {
	lat := chi.URLParam(r, "lat")
	long := chi.URLParam(r, "long")

	if lat == "" || long == "" {
		return 0, 0, fmt.Errorf("latitude and longitude must be provided")
	}

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse latitude: %w", err)
	}

	longFloat, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return latFloat, longFloat, nil
}
