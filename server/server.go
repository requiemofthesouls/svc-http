package server

import (
	"errors"
	"net/http"
)

type (
	Server interface {
		Start() error
		Stop() error
		IsStarted() bool
	}

	server struct {
		server    *http.Server
		isStarted bool
	}

	Config struct {
		Name    string `mapstructure:"name"`
		Address string `mapstructure:"address"`
	}
)

func New(config Config, handler Handler) Server {
	return &server{
		server: &http.Server{
			Addr:    config.Address,
			Handler: handler,
		},
	}
}

func (s *server) Start() error {
	s.isStarted = true
	defer func() { s.isStarted = false }()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *server) Stop() error {
	return s.server.Close()
}

func (s *server) IsStarted() bool {
	return s.isStarted
}
