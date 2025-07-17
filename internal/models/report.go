// internal/models/report.go
package models

import "time"

type Report struct {
	ID              int       `json:"id"`
	MetricName      string    `json:"metric_name"`
	PatternDetected string    `json:"pattern_detected"`
	Severity        string    `json:"severity"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Details         string    `json:"details"`
	CreatedAt       time.Time `json:"created_at"`
}