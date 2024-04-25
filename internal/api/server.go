package api

import (
	"construction-organization-report/internal/log"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	listenAddr string
	router     *mux.Router
}

func NewServer(listenAddr string, router *mux.Router) *Server {
	return &Server{listenAddr: listenAddr, router: router}
}

func (s *Server) InitializeRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/report", handleGetLastProjectsReports).Methods("GET")
	api.HandleFunc("/report/{id:[0-9]+}/create", handleCreateReport).Methods("GET")
	api.HandleFunc("/report/{id:[0-9]+}", handleGetProjectReports).Methods("GET")
}

func (s *Server) Start() {
	log.Logger.Fatal(http.ListenAndServe(s.listenAddr, s.router))
}
