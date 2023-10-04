package weatherservice_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-sample-rest/internal/types"
	"go-sample-rest/internal/weatherservice"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type MockWeatherDataClient struct {
	getLatestWeatherData func(lat, long float64) (*types.WeatherData, error)
}

func (m *MockWeatherDataClient) GetLatestWeatherData(lat, long float64) (*types.WeatherData, error) {
	if m != nil && m.getLatestWeatherData != nil {
		return m.getLatestWeatherData(lat, long)
	}

	return &types.WeatherData{
		Latitude:      lat,
		Longitude:     long,
		Temperature:   3.3,
		WindSpeed:     4.4,
		WindDirection: 5.5,
		CreatedAt:     "2023-10-04T06:53:38.581587Z",
	}, nil
}

type MockWeatherDataRepository struct {
	getLatestWeatherData func(lat, long float64) (*types.WeatherData, error)
	getWeatherHistory    func(lat, long float64) ([]*types.WeatherData, error)
	saveWeatherData      func(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error)
}

func (m *MockWeatherDataRepository) GetLatestWeatherData(lat, long float64) (*types.WeatherData, error) {
	if m != nil && m.getLatestWeatherData != nil {
		return m.getLatestWeatherData(lat, long)
	}

	return &types.WeatherData{
		Latitude:      lat,
		Longitude:     long,
		Temperature:   3.3,
		WindSpeed:     4.4,
		WindDirection: 5.5,
		CreatedAt:     "2023-10-04T06:53:38.581587Z",
	}, nil
}

func (m *MockWeatherDataRepository) GetWeatherHistory(lat, long float64) ([]*types.WeatherData, error) {
	if m != nil && m.getWeatherHistory != nil {
		return m.getWeatherHistory(lat, long)
	}

	return []*types.WeatherData{
		{
			Latitude:      lat,
			Longitude:     long,
			Temperature:   3.3,
			WindSpeed:     4.4,
			WindDirection: 5.5,
			CreatedAt:     "2023-10-04T06:53:38.581587Z",
		},
	}, nil
}

func (m *MockWeatherDataRepository) SaveWeatherData(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error) {
	if m != nil && m.saveWeatherData != nil {
		return m.saveWeatherData(lat, long, temperature, windDirection, windSpeed)
	}

	return &types.WeatherData{
		Latitude:      lat,
		Longitude:     long,
		Temperature:   3.3,
		WindSpeed:     4.4,
		WindDirection: 5.5,
		CreatedAt:     "2023-10-04T06:53:38.581587Z",
	}, nil
}

func TestGetLatestWeatherData(t *testing.T) {
	testCases := []struct {
		name                      string
		lat                       string
		long                      string
		mockWeatherDataRepository *MockWeatherDataRepository
		expectedStatusCode        int
		expectedBody              string
	}{
		{
			name:                      "should err when no lat and long provided",
			lat:                       "",
			long:                      "",
			mockWeatherDataRepository: &MockWeatherDataRepository{},
			expectedStatusCode:        http.StatusBadRequest,
			expectedBody:              http.StatusText(http.StatusBadRequest),
		},
		{
			name: "should return internal error when repo returns an error trying to get latest weather data",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return nil, fmt.Errorf("error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name: "should return not found when repo did not return an error but weatherData is nil",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return nil, nil
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       http.StatusText(http.StatusNotFound),
		},
		{
			name: "should return weather data as json when repo returns weather data",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return &types.WeatherData{
						Id:            "abc123",
						Latitude:      lat,
						Longitude:     long,
						Temperature:   3.3,
						WindSpeed:     4.4,
						WindDirection: 5.5,
						CreatedAt:     "2023-10-04T06:53:38.581587Z",
					}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":"abc123","latitude":1.1,"longitude":2.2,"temperature":3.3,"wind_direction":5.5,"wind_speed":4.4,"created_at":"2023-10-04T06:53:38.581587Z"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := weatherservice.NewService(&MockWeatherDataClient{}, tc.mockWeatherDataRepository)
			r := httptest.NewRequest("GET", fmt.Sprintf("/%s,%s/latest", tc.lat, tc.long), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("lat", tc.lat)
			rctx.URLParams.Add("long", tc.long)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			service.GetLatestWeather(w, r)

			require.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
			require.Equal(t, tc.expectedBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestGetWeatherHistory(t *testing.T) {
	testCases := []struct {
		name                      string
		lat                       string
		long                      string
		mockWeatherDataRepository *MockWeatherDataRepository
		expectedStatusCode        int
		expectedBody              string
	}{
		{
			name:                      "should err when no lat and long provided",
			lat:                       "",
			long:                      "",
			mockWeatherDataRepository: &MockWeatherDataRepository{},
			expectedStatusCode:        http.StatusBadRequest,
			expectedBody:              http.StatusText(http.StatusBadRequest),
		},
		{
			name: "should return internal error when repo returns an error trying to get weather data history",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				getWeatherHistory: func(lat, long float64) ([]*types.WeatherData, error) {
					return nil, fmt.Errorf("error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name: "should return weather data history as json when repo returns weather data history",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				getWeatherHistory: func(lat, long float64) ([]*types.WeatherData, error) {
					return []*types.WeatherData{
						{
							Id:            "abc123",
							Latitude:      lat,
							Longitude:     long,
							Temperature:   3.3,
							WindSpeed:     4.4,
							WindDirection: 5.5,
							CreatedAt:     "2023-10-04T06:53:38.581587Z",
						},
					}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":"abc123","latitude":1.1,"longitude":2.2,"temperature":3.3,"wind_direction":5.5,"wind_speed":4.4,"created_at":"2023-10-04T06:53:38.581587Z"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := weatherservice.NewService(&MockWeatherDataClient{}, tc.mockWeatherDataRepository)
			r := httptest.NewRequest("GET", fmt.Sprintf("/%s,%s/history", tc.lat, tc.long), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("lat", tc.lat)
			rctx.URLParams.Add("long", tc.long)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			service.GetWeatherHistory(w, r)

			require.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
			require.Equal(t, tc.expectedBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}

func TestUpdateWeather(t *testing.T) {
	testCases := []struct {
		name                      string
		lat                       string
		long                      string
		mockWeatherDataRepository *MockWeatherDataRepository
		mockWeatherDataClient     *MockWeatherDataClient
		expectedStatusCode        int
		expectedBody              string
	}{
		{
			name:                      "should err when no lat and long provided",
			lat:                       "",
			long:                      "",
			mockWeatherDataRepository: &MockWeatherDataRepository{},
			mockWeatherDataClient:     &MockWeatherDataClient{},
			expectedStatusCode:        http.StatusBadRequest,
			expectedBody:              http.StatusText(http.StatusBadRequest),
		},
		{
			name:                      "should return internal error when data client returns an error trying to get latest weather data",
			lat:                       "1.1",
			long:                      "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{},
			mockWeatherDataClient: &MockWeatherDataClient{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return nil, fmt.Errorf("error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name:                      "should return status not found when data client did not return an error but weatherData is nil",
			lat:                       "1.1",
			long:                      "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{},
			mockWeatherDataClient: &MockWeatherDataClient{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return nil, nil
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       http.StatusText(http.StatusNotFound),
		},
		{
			name: "should return internal error when repo returns an error trying to save weather data",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				saveWeatherData: func(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error) {
					return nil, fmt.Errorf("error")
				},
			},
			mockWeatherDataClient: &MockWeatherDataClient{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return &types.WeatherData{
						Latitude:      lat,
						Longitude:     long,
						Temperature:   3.3,
						WindSpeed:     4.4,
						WindDirection: 5.5,
					}, nil
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       http.StatusText(http.StatusInternalServerError),
		},
		{
			name: "should return OK with weather data as json when repo returns weather data",
			lat:  "1.1",
			long: "2.2",
			mockWeatherDataRepository: &MockWeatherDataRepository{
				saveWeatherData: func(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error) {
					return &types.WeatherData{
						Id:            "abc123",
						Latitude:      lat,
						Longitude:     long,
						Temperature:   temperature,
						WindSpeed:     windDirection,
						WindDirection: windSpeed,
						CreatedAt:     "2023-10-04T06:53:38.581587Z",
					}, nil
				},
			},
			mockWeatherDataClient: &MockWeatherDataClient{
				getLatestWeatherData: func(lat, long float64) (*types.WeatherData, error) {
					return &types.WeatherData{
						Latitude:      lat,
						Longitude:     long,
						Temperature:   3.3,
						WindSpeed:     4.4,
						WindDirection: 5.5,
					}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":"abc123","latitude":1.1,"longitude":2.2,"temperature":3.3,"wind_direction":4.4,"wind_speed":5.5,"created_at":"2023-10-04T06:53:38.581587Z"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := weatherservice.NewService(tc.mockWeatherDataClient, tc.mockWeatherDataRepository)
			r := httptest.NewRequest("POST", fmt.Sprintf("/%s,%s/update", tc.lat, tc.long), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("lat", tc.lat)
			rctx.URLParams.Add("long", tc.long)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			service.UpdateWeather(w, r)

			require.Equal(t, tc.expectedStatusCode, w.Result().StatusCode)
			require.Equal(t, tc.expectedBody, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
