package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/config"
	"github.com/VladimirStepanov/urlshortener/pkg/shortener"
	"github.com/VladimirStepanov/urlshortener/pkg/store"
	"github.com/sirupsen/logrus"
)

//Server ...
type Server struct {
	log       *logrus.Logger
	db        store.Storage
	config    *config.Config
	shortener shortener.Shortener
}

func getLogger(level string) (*logrus.Logger, error) {
	if level == "" {
		level = "INFO"
	}

	lvl, err := logrus.ParseLevel(level)

	if err != nil {
		return nil, err
	}

	log := logrus.New()

	log.SetLevel(lvl)

	return log, nil
}

//New ...
func New(cfg *config.Config, dbConn store.Storage, shortener shortener.Shortener) (*Server, error) {
	log, err := getLogger(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	return &Server{log: log, db: dbConn, config: cfg, shortener: shortener}, nil
}

//Start run server
func (s *Server) Start() error {

	srv := &http.Server{
		Handler:      s.router(),
		Addr:         fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.log.Infof("Starting server on %s:%s\n", s.config.Host, s.config.Port)

	return srv.ListenAndServe()
}
