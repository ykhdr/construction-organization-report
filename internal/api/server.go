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

func NewServer(listenAddr string, router *mux.Router, db *sql.DB) *Server {
	return &Server{listenAddr: listenAddr, router: router, db: db}
}

func (s *Server) InitializeRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/report", s.handleGetLastProjectsReports).Methods("GET")
	api.HandleFunc("/report/{id:[0-9]+}/create", s.handleCreateReport).Methods("GET")
	api.HandleFunc("/report/{id:[0-9]+}", s.handleGetProjectReports).Methods("GET")
}

func (s *Server) Start() {
	log.Logger.Fatal(http.ListenAndServe(s.listenAddr, s.router))
}
