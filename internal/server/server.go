package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/berberapan/dota-work/internal/handlers"
)

type Server struct {
	httpAddr  string
	httpsAddr string
	router    *http.ServeMux
	db        *sql.DB
}

func NewServer(httpAddr, httpsAddr string, router *http.ServeMux, db *sql.DB) *Server {
	return &Server{
		httpAddr:  httpAddr,
		httpsAddr: httpsAddr,
		router:    router,
		db:        db,
	}
}

func (s *Server) setupRoutes() {
	subrouter := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web/dist"))
	s.router.Handle("/app/", http.StripPrefix("/app", fs))

	s.router.Handle("/v1/", http.StripPrefix("/v1", subrouter))

	subrouter.HandleFunc("GET /healthz", handlers.HealthCheck)
	subrouter.HandleFunc("POST /teamdata", handlers.GetTeamData)
	subrouter.HandleFunc("POST /tournamentschedule", handlers.GetTournamentSchedule)
}

func (s *Server) Run() error {
	s.setupRoutes()

	go func() {
		log.Printf("HTTP server starting on  %s (redirecting to HTTPS)", s.httpAddr)
		if err := http.ListenAndServe(s.httpAddr, http.HandlerFunc(s.redirectTLS)); err != nil {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
	}()

	log.Printf("HTTPS server starting on %s", s.httpsAddr)
	return http.ListenAndServeTLS(s.httpsAddr,
		"/etc/letsencrypt/live/dota.boukdir.se/fullchain.pem",
		"/etc/letsencrypt/live/dota.boukdir.se/privkey.pem",
		s.router)
}

func (s *Server) redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}
