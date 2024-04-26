package api

import (
	"construction-organization-report/internal/log"
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
		log.Logger.Warningln("Error in parse project id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Logger.Infoln("Create report for project", projectID)
	reportFile, err := report.CreateReport(projectID, s.db)
	if err != nil {
		log.Logger.Warningln("Error on create report", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Logger.Infoln("Report created", projectID)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reportFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) handleGetReports(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	projectID, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Logger.Warningln("Error in parse project id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Logger.Infoln("Get report for project", projectID)
	rawReport, err := report.GetReports(projectID, s.db)
	if err != nil {
		log.Logger.Warningln("Error on get report", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Logger.Infoln("Report got successful", projectID)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rawReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
