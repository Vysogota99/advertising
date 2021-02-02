package server

import (
	"github.com/Vysogota99/advertising/internal/app/store/postgres"

	"github.com/Vysogota99/advertising/internal/app/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //
)

// Server ...
type Server struct {
	Conf   *Config
	router *Router
	store  store.Store
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	return &Server{
		Conf: conf,
	}, nil
}

// Start - start the server
func (s *Server) Start() error {
	if err := s.initStore(); err != nil {
		return err
	}

	s.initRouter()

	s.router.Setup().Run(s.Conf.serverPort)
	return nil
}

func (s *Server) initRouter() {
	router := NewRouter(s.Conf.serverPort, s.store)
	s.router = router
}

func (s *Server) initStore() error {
	db, err := sqlx.Connect("postgres", s.Conf.dbConnString)
	if err != nil {
		return err
	}

	s.store = postgres.New(db)

	return nil
}
