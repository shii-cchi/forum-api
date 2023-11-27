package server

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shii-cchi/forum-api/internal/database"
	"github.com/shii-cchi/forum-api/internal/handlers"
	"github.com/shii-cchi/forum-api/pkg/hash"
	"log"
	"net/http"
	"os"
)

type Server struct {
	httpServer  *http.Server
	httpHandler *handlers.Handler
	queries     *database.Queries
}

func NewServer(r chi.Router) (*Server, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("PORT is not found")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, errors.New("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	salt := os.Getenv("SALT_STRING")

	if salt == "" {
		return nil, errors.New("SALT_STRING is not found")
	}

	queries := database.New(conn)
	hasher := hash.NewSHA1Hasher(salt)

	handler := handlers.New(queries, hasher)
	handler.RegisterHTTPEndpoints(r)

	log.Printf("Server starting on port %s", port)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
		httpHandler: handler,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
