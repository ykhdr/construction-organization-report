package api

import (
	"construction-organization-report/internal/log"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	listenAddr string
	router     *mux.Router
	db         *sql.DB
}

func NewServer(listenAddr string, db *sql.DB) *Server {
	return &Server{listenAddr: listenAddr, router: &mux.Router{}, db: db}
}

func (s *Server) InitializeRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/report/create", s.handleCreateReport).Methods("GET")
	api.HandleFunc("/report", s.handleGetReports).Methods("GET")
	api.HandleFunc("/report/{id:[0-9]+}", s.handleGetReport).Methods("GET")
}

func (s *Server) Start() {
	log.Logger.Fatal(http.ListenAndServe(s.listenAddr, s.router))
}
