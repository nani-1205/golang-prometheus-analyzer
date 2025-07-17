// internal/api/handlers.go
package api

import (
	"encoding/json"
	"html/template"
	//"log"//
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/analysis"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/config"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/database"
)

var templates = template.Must(template.ParseFiles("web/templates/index.html"))

func RegisterHandlers(r *mux.Router, cfg *config.Config) {
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/reports", getReportsHandler).Methods("GET")
	r.HandleFunc("/api/analyze/cpu", func(w http.ResponseWriter, r *http.Request) {
		analyzeCPUHandler(w, r, cfg)
	}).Methods("POST")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getReportsHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := database.GetAllReports()
	if err != nil {
		http.Error(w, "Failed to fetch reports", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func analyzeCPUHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	// Run analysis in a goroutine so the API call returns immediately
	go analysis.AnalyzeCPUUsagePattern(cfg)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "CPU analysis started in the background."})
}