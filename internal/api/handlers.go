// internal/api/handlers.go
package api

import (
	"encoding/json"
	"html/template"
	"log"
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

	// --- UPDATED ENDPOINTS ---

	// Endpoint for the original "transient spike" analysis
	r.HandleFunc("/api/analyze/cpu-spike", func(w http.ResponseWriter, r *http.Request) {
		go analysis.AnalyzeCPUUsageSpikePattern(cfg)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"message": "CPU spike analysis started."})
	}).Methods("POST")

	// Endpoint for the NEW "sustained high load" analysis
	r.HandleFunc("/api/analyze/cpu-high-load", func(w http.ResponseWriter, r *http.Request) {
		go analysis.AnalyzeSustainedHighCPU(cfg)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sustained high CPU analysis started."})
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
		log.Printf("Error fetching reports: %v", err)
		http.Error(w, "Failed to fetch reports", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}