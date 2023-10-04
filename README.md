# go-sample-rest

Sample Go REST application.

This is a basic REST API that checks the weather conditions of a specific location. It queries it from open-mateo, and store it in a persistent DB.

# Endpoints

## GET /weather/{lat},{long}/latest
This endpoint gets the latest weather information stored in DB

## GET /weather/{lat},{long}/history
This endpoint pulls all weather information stored in DB

## POST /weather/{lat},{long}/update
This endpoint pulls latest weather information from OpenMateo, adds another entry in DB which then becomes the latest weather data for this location

# Local usage

Start local container depdendencies
```shell script
docker-compose up
```

Start local development server
```shell script
make dev
```

# Tests

```shell script
# Running unit tests
make tests

# Running integration tests
docker-compose -f docker-compose.integration.yml up # This fires up integration postgres instance
make integration_tests
```
