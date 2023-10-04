package main

import (
	"database/sql"
	"net/http"

	"go-sample-rest/internal/openmateo"
	"go-sample-rest/internal/repository"
	"go-sample-rest/internal/server"
	"go-sample-rest/internal/weatherservice"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialise config
	config := NewConfig()

	// Initialise repository
	db, err := sql.Open("postgres", config.PGConnString)
	if err != nil {
		log.Fatalf("failed to initialise db: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	// Initialise open mateo client
	httpClient := &http.Client{
		Timeout: config.HTTPTimeout,
	}
	openMateoClient := openmateo.NewClient(httpClient)

	// Initialise weather service
	weatherService := weatherservice.NewService(openMateoClient, repo)

	s := server.NewServer(config.Port, weatherService)
	s.Start()
}
