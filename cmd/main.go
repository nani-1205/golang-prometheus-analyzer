// cmd/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/your-username/golang-prometheus-analyzer/internal/api"
	"github.com/your-username/golang-prometheus-analyzer/internal/config"
	"github.com/your-username/golang-prometheus-analyzer/internal/database"
	"github.com/your-username/golang-prometheus-analyzer/internal/reporting"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	// 2. Initialize Database
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer database.CloseDB()
	log.Println("Database connection successful.")

	// 3. Start the daily report scheduler
	reporting.StartScheduler(cfg)
	log.Println("Reporting scheduler started.")

	// 4. Setup Router and Handlers
	r := mux.NewRouter()
	api.RegisterHandlers(r, cfg)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Printf("Starting server on port %s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}