package server

import (
	"database/sql"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/shii-cchi/forum-api/internal/config"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers"
	"github.com/shii-cchi/forum-api/internal/service"
	"github.com/shii-cchi/forum-api/pkg/hash"
	"log"
	"net/http"
)

type Server struct {
	httpServer  *http.Server
	httpHandler *handlers.Handler
	queries     *database.Queries
}

func NewServer(r chi.Router) (*Server, error) {
	cfg, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("postgres", cfg.DbURI)

	if err != nil {
		return nil, err
	}

	queries := database.New(conn)
	hasher := hash.NewSHA1Hasher(cfg.SaltString)

	services := service.NewServices(service.Deps{
		Queries: queries,
		Hasher:  hasher,
		Config:  cfg,
	})

	handler := handlers.NewHandler(services, cfg)
	handler.RegisterHTTPEndpoints(r)

	log.Printf("Server starting on port %s", cfg.Port)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: r,
		},
		httpHandler: handler,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
