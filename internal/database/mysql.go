// internal/database/mysql.go
package database

import (
	"database/sql"
	"fmt"
	"time" // Make sure 'time' is imported

	_ "github.com/go-sql-driver/mysql"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/config"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/models"
)

var db *sql.DB

func InitDB(cfg *config.Config) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	return db.Ping()
}

func CloseDB() {
	db.Close()
}

func InsertReport(report models.Report) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO reports(metric_name, pattern_detected, severity, start_time, end_time, details) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(report.MetricName, report.PatternDetected, report.Severity, report.StartTime, report.EndTime, report.Details)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetAllReports() ([]models.Report, error) {
	rows, err := db.Query("SELECT id, metric_name, pattern_detected, severity, start_time, end_time, details, created_at FROM reports ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var r models.Report
		if err := rows.Scan(&r.ID, &r.MetricName, &r.PatternDetected, &r.Severity, &r.StartTime, &r.EndTime, &r.Details, &r.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

// --- NEW FUNCTION ---
// CheckForRecentReport checks if a report for a specific metric/instance exists within the cooldown period.
func CheckForRecentReport(metricName string, instance string, cooldown time.Duration) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM reports 
			WHERE metric_name = ? 
			AND details LIKE ? 
			AND created_at >= ?
		)`

	// We use LIKE to find the instance name within the details string
	instancePattern := "%" + instance + "%"

	// Calculate the time threshold
	sinceTime := time.Now().Add(-cooldown)

	err := db.QueryRow(query, metricName, instancePattern, sinceTime).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}