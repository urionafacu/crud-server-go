package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urionafacu/api-rest-websockets-go/database"
	"github.com/urionafacu/api-rest-websockets-go/repository"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(_ context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is missing")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWT secret is missing")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("database URL is missing")
	}

	return &Broker{
		config: config,
		router: mux.NewRouter(),
	}, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)

	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)

	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
