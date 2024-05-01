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
	projectID, err := strconv.Atoi(r.URL.Query().Get("project_id"))
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
	projectID, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		log.Logger.Warningln("Error in parse project id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Logger.Infoln("Get report for project", projectID)
	rawReport, err := report.GetRawReports(projectID, s.db)
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

func (s *Server) handleGetReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reportID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Logger.Warningln("Error in parse report id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Logger.Infoln("Get report", reportID)
	rawReport, err := report.GetRawReport(reportID, s.db)
	if err != nil {
		log.Logger.Warningln("Error on get report", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Logger.Infoln("Report got successful", reportID)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rawReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
