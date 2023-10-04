//go:build integration

package repository_test

import (
	"database/sql"
	"go-sample-rest/internal/repository"
	"go-sample-rest/internal/types"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
)

var pgConnString = "postgres://weather:weather@localhost:7432/weather?sslmode=disable"

func TestIntegrationGetLatestWeatherData(t *testing.T) {
	testCases := []struct {
		name                string
		setup               func(db *sql.DB, t *testing.T)
		shouldError         bool
		expectedWeatherData *types.WeatherData
	}{
		{
			name: "should not error if weather data does not exist",
			setup: func(db *sql.DB, t *testing.T) {
				// No setup required
			},
			shouldError:         false,
			expectedWeatherData: nil,
		},
		{
			name: "should return weather data if it exists in db",
			setup: func(db *sql.DB, t *testing.T) {
				_, err := db.Exec(`
          INSERT INTO "weather"."weather_data" (id, latitude, longitude, temperature, wind_speed, wind_direction, created_at)
          VALUES ('abc123', 1.1, 2.2, 3.3, 4.4, 5.5, '2023-10-04T06:53:38.581587Z')
        `)

				require.NoError(t, err)
			},
			shouldError: false,
			expectedWeatherData: &types.WeatherData{
				Id:            "abc123",
				Latitude:      1.1,
				Longitude:     2.2,
				Temperature:   3.3,
				WindSpeed:     4.4,
				WindDirection: 5.5,
				CreatedAt:     "2023-10-04T06:53:38.581587Z",
			},
		},
	}

	// Initialise db connection
	dbClient, err := sql.Open("postgres", pgConnString)
	if err != nil {
		log.Fatalf("failed to initialise db: %v", err)
	}
	defer dbClient.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(dbClient, t)

			repository := repository.NewRepository(dbClient)
			weatherData, err := repository.GetLatestWeatherData(1.1, 2.2)

			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedWeatherData, weatherData)
			}
		})

		_, err = dbClient.Exec(`DELETE FROM "weather"."weather_data"`)
		require.NoError(t, err)
	}
}

func TestIntegrationGetWeatherHistory(t *testing.T) {
	testCases := []struct {
		name                   string
		setup                  func(db *sql.DB, t *testing.T)
		shouldError            bool
		expectedWeatherHistory []*types.WeatherData
	}{
		{
			name: "should return weather history if it exists in db",
			setup: func(db *sql.DB, t *testing.T) {
				_, err := db.Exec(`
          INSERT INTO "weather"."weather_data" (id, latitude, longitude, temperature, wind_speed, wind_direction, created_at)
          VALUES ('a1', 1.1, 2.2, 1.3, 1.4, 1.5, '2023-10-04T06:53:38.581587Z');

          INSERT INTO "weather"."weather_data" (id, latitude, longitude, temperature, wind_speed, wind_direction, created_at)
          VALUES ('a2', 1.1, 2.2, 2.3, 2.4, 2.5, '2023-10-04T06:55:38.581587Z');

          INSERT INTO "weather"."weather_data" (id, latitude, longitude, temperature, wind_speed, wind_direction, created_at)
          VALUES ('a3', 1.1, 2.2, 3.3, 3.4, 3.5, '2023-10-04T06:57:38.581587Z');
        `)

				require.NoError(t, err)
			},
			shouldError: false,
			expectedWeatherHistory: []*types.WeatherData{
				{
					Id:            "a3",
					Latitude:      1.1,
					Longitude:     2.2,
					Temperature:   3.3,
					WindSpeed:     3.4,
					WindDirection: 3.5,
					CreatedAt:     "2023-10-04T06:57:38.581587Z",
				},
				{
					Id:            "a2",
					Latitude:      1.1,
					Longitude:     2.2,
					Temperature:   2.3,
					WindSpeed:     2.4,
					WindDirection: 2.5,
					CreatedAt:     "2023-10-04T06:55:38.581587Z",
				},
				{
					Id:            "a1",
					Latitude:      1.1,
					Longitude:     2.2,
					Temperature:   1.3,
					WindSpeed:     1.4,
					WindDirection: 1.5,
					CreatedAt:     "2023-10-04T06:53:38.581587Z",
				},
			},
		},
	}

	// Initialise db connection
	dbClient, err := sql.Open("postgres", pgConnString)
	if err != nil {
		log.Fatalf("failed to initialise db: %v", err)
	}
	defer dbClient.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(dbClient, t)

			repository := repository.NewRepository(dbClient)
			weatherData, err := repository.GetWeatherHistory(1.1, 2.2)

			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedWeatherHistory, weatherData)
			}
		})

		_, err = dbClient.Exec(`DELETE FROM "weather"."weather_data"`)
		require.NoError(t, err)
	}
}

func TestIntegrationSaveWeatherData(t *testing.T) {
	testCases := []struct {
		name                string
		shouldError         bool
		lat                 float64
		long                float64
		temperature         float64
		windSpeed           float64
		windDirection       float64
		assert              func(repo *repository.Repository, savedWeatherData *types.WeatherData, t *testing.T)
		expectedWeatherData *types.WeatherData
	}{
		{
			name:          "should return weather history if it exists in db",
			shouldError:   false,
			lat:           1.1,
			long:          2.2,
			temperature:   3.3,
			windSpeed:     4.4,
			windDirection: 5.5,
			assert: func(repo *repository.Repository, savedWeatherData *types.WeatherData, t *testing.T) {
				weatherData, err := repo.GetLatestWeatherData(1.1, 2.2)

				require.Equal(t, weatherData, savedWeatherData)
				require.NoError(t, err)
				require.Equal(t, &types.WeatherData{
					Id:            weatherData.Id,
					Latitude:      1.1,
					Longitude:     2.2,
					Temperature:   3.3,
					WindSpeed:     4.4,
					WindDirection: 5.5,
					CreatedAt:     weatherData.CreatedAt,
				}, weatherData)
			},
		},
	}

	// Initialise db connection
	dbClient, err := sql.Open("postgres", pgConnString)
	if err != nil {
		log.Fatalf("failed to initialise db: %v", err)
	}
	defer dbClient.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := repository.NewRepository(dbClient)
			savedWeatherData, err := repository.SaveWeatherData(tc.lat, tc.long, tc.temperature, tc.windDirection, tc.windSpeed)
			require.NoError(t, err)

			tc.assert(repository, savedWeatherData, t)
		})

		_, err = dbClient.Exec(`DELETE FROM "weather"."weather_data"`)
		require.NoError(t, err)
	}
}
