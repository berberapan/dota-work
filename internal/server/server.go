package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/berberapan/dota-work/internal/handlers"
)

type Server struct {
	listenAddr string
	router     *http.ServeMux
	db         *sql.DB
}

func NewServer(listenAddr string, router *http.ServeMux, db *sql.DB) *Server {
	return &Server{
		listenAddr: listenAddr,
		router:     router,
		db:         db,
	}
}

func (s *Server) setupRoutes() {
	subrouter := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web/dist"))
	s.router.Handle("/app/", http.StripPrefix("/app", fs))

	s.router.Handle("/v1/", http.StripPrefix("/v1", subrouter))

	subrouter.HandleFunc("GET /healthz", handlers.HealthCheck)
}

func (s *Server) Run() error {
	s.setupRoutes()
	log.Printf("Server starting and listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}
