package openmateo_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-sample-rest/internal/openmateo"
	"go-sample-rest/internal/types"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var mockResponseBody = openmateo.OpenMateoForecastResponseBody{
	Latitude:  1.1,
	Longitude: 2.2,
	CurrentWeather: openmateo.OpenMateoCurrentWeather{
		Temperature:   3.3,
		WindSpeed:     4.4,
		WindDirection: 5.5,
	},
}

type MockHTTPClient struct {
	get func(url string) (resp *http.Response, err error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if m != nil && m.get != nil {
		return m.get(url)
	}

	jsonBody, _ := json.Marshal(mockResponseBody)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(jsonBody)),
	}, nil
}

func TestGetLatestWeatherData(t *testing.T) {
	testCases := []struct {
		name             string
		mockHTTPClient   *MockHTTPClient
		shouldError      bool
		expectedResponse *types.WeatherData
	}{
		{
			name: "should return error when http client returns error",
			mockHTTPClient: &MockHTTPClient{
				get: func(url string) (resp *http.Response, err error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBuffer([]byte("some error"))),
					}, fmt.Errorf("some error")
				},
			},
			shouldError: true,
		},
		{
			name: "should return error when http client returns non-200 status code",
			mockHTTPClient: &MockHTTPClient{
				get: func(url string) (resp *http.Response, err error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBuffer([]byte("some error"))),
					}, nil
				},
			},
			shouldError: true,
		},
		{
			name: "should return error when http client returns invalid json",
			mockHTTPClient: &MockHTTPClient{
				get: func(url string) (resp *http.Response, err error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer([]byte("invalid json"))),
					}, nil
				},
			},
			shouldError: true,
		},
		{
			name: "should return weather data if there is no error",
			mockHTTPClient: &MockHTTPClient{
				get: func(url string) (resp *http.Response, err error) {
					jsonBody, _ := json.Marshal(mockResponseBody)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(jsonBody)),
					}, nil
				},
			},
			shouldError: false,
			expectedResponse: &types.WeatherData{
				Latitude:      1.1,
				Longitude:     2.2,
				Temperature:   3.3,
				WindSpeed:     4.4,
				WindDirection: 5.5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			openMateo := openmateo.NewClient(tc.mockHTTPClient)

			weatherData, err := openMateo.GetLatestWeatherData(1.1, 2.2)

			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.Equal(t, tc.expectedResponse, weatherData)
			}
		})
	}
}
