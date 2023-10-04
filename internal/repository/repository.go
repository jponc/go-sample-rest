package repository

import (
	"database/sql"

	"go-sample-rest/internal/types"

	"github.com/google/uuid"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) *Repository {
	return &Repository{
		dbClient: dbClient,
	}
}

func (r *Repository) GetLatestWeatherData(lat, long float64) (*types.WeatherData, error) {
	row := r.dbClient.QueryRow(`
    SELECT id, latitude, longitude, temperature, wind_direction, wind_speed, created_at
    FROM weather_data
    WHERE latitude = $1 AND longitude = $2
    ORDER BY created_at DESC
    LIMIT 1
  `, lat, long)

	var weatherData types.WeatherData

	err := row.Scan(
		&weatherData.Id,
		&weatherData.Latitude,
		&weatherData.Longitude,
		&weatherData.Temperature,
		&weatherData.WindDirection,
		&weatherData.WindSpeed,
		&weatherData.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No data found
			return nil, nil
		}
		return nil, err
	}

	return &weatherData, nil
}

func (r *Repository) GetWeatherHistory(lat, long float64) ([]*types.WeatherData, error) {
	rows, err := r.dbClient.Query(`
    SELECT id, latitude, longitude, temperature, wind_direction, wind_speed, created_at
    FROM weather_data
    WHERE latitude = $1 AND longitude = $2
    ORDER BY created_at DESC
  `, lat, long)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var weatherDataList []*types.WeatherData

	for rows.Next() {
		var weatherData types.WeatherData

		err := rows.Scan(
			&weatherData.Id,
			&weatherData.Latitude,
			&weatherData.Longitude,
			&weatherData.Temperature,
			&weatherData.WindDirection,
			&weatherData.WindSpeed,
			&weatherData.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		weatherDataList = append(weatherDataList, &weatherData)
	}

	return weatherDataList, nil
}

func (r *Repository) SaveWeatherData(lat, long, temperature, windDirection, windSpeed float64) (*types.WeatherData, error) {
	id := uuid.New()

	row := r.dbClient.QueryRow(`
    INSERT INTO weather_data (id, latitude, longitude, temperature, wind_direction, wind_speed)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, latitude, longitude, temperature, wind_direction, wind_speed, created_at
  `, id.String(), lat, long, temperature, windDirection, windSpeed)

	var weatherData types.WeatherData

	err := row.Scan(
		&weatherData.Id,
		&weatherData.Latitude,
		&weatherData.Longitude,
		&weatherData.Temperature,
		&weatherData.WindDirection,
		&weatherData.WindSpeed,
		&weatherData.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &weatherData, nil
}
