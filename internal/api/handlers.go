package api

import (
	"construction-organization-report/internal/report"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) handleCreateReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reportFile, err := report.CreateReport(projectID, s.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reportFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) handleGetLastProjectsReports(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleGetProjectReports(w http.ResponseWriter, r *http.Request) {

}
