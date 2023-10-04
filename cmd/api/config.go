package main

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	PGConnString string
	HTTPTimeout  time.Duration
	Port         int
}

func NewConfig() *Config {
	port, err := strconv.Atoi(getEnv("PORT"))
	if err != nil {
		log.Fatalf("Cannot convert port to int")
	}

	return &Config{
		PGConnString: getEnv("PG_DB_CONN_STRING"),
		HTTPTimeout:  time.Second * 10,
		Port:         port,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("Environment variable %s not found", key)
	}

	return val
}
