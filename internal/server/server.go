package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	port           int
	weatherService WeatherService
}

type WeatherService interface {
	GetLatestWeather(w http.ResponseWriter, r *http.Request)
	GetWeatherHistory(w http.ResponseWriter, r *http.Request)
	UpdateWeather(w http.ResponseWriter, r *http.Request)
}

func NewServer(port int, weatherService WeatherService) *Server {
	return &Server{
		port:           port,
		weatherService: weatherService,
	}
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/weather", func(r chi.Router) {
		r.Get("/{lat},{long}/latest", s.weatherService.GetLatestWeather)
		r.Get("/{lat},{long}/history", s.weatherService.GetWeatherHistory)
		r.Post("/{lat},{long}/update", s.weatherService.UpdateWeather)
	})

	log.Infof("Starting server on port %d", s.port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}
