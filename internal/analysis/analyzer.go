// internal/analysis/analyzer.go
package analysis

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/common/model"
	"github.com/your-username/golang-prometheus-analyzer/internal/config"
	"github.com/your-username/golang-prometheus-analyzer/internal/database"
	"github.com/your-username/golang-prometheus-analyzer/internal/models"
	"github.com/your-username/golang-prometheus-analyzer/internal/notification"
	prom "github.com/your-username/golang-prometheus-analyzer/internal/prometheus"
)

// AnalyzeCPUUsagePattern is our example analyzer.
// It checks for a sudden spike in CPU usage in a specific time window.
func AnalyzeCPUUsagePattern(cfg *config.Config) {
	log.Println("Running CPU usage pattern analysis...")
	now := time.Now()
	// Let's analyze the last hour's data.
	start := now.Add(-1 * time.Hour)
	end := now

	// PromQL query to get average CPU usage (non-idle time)
	// This gives a value between 0-100 per instance
	query := `100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[2m])) * 100)`

	result, err := prom.QueryRange(cfg, query, start, end, 1*time.Minute)
	if err != nil {
		log.Printf("Error querying Prometheus: %v", err)
		return
	}

	matrix, ok := result.(model.Matrix)
	if !ok {
		log.Printf("Error: Expected a matrix result type, got %T", result)
		return
	}

	for _, stream := range matrix {
		instance := string(stream.Metric["instance"])
		checkAndReportSpike(cfg, stream.Values, instance)
	}
	log.Println("CPU usage pattern analysis finished.")
}

// A simple algorithm to detect a spike
func checkAndReportSpike(cfg *config.Config, values []model.SamplePair, instance string) {
	const spikeThreshold = 70.0 // CPU %
	const baseThreshold = 20.0  // Normal CPU %
	const jumpPercentage = 50.0 // The jump must be at least 50%

	var spikeStartTime, spikeEndTime time.Time
	inSpike := false

	for i := 1; i < len(values); i++ {
		prevValue := float64(values[i-1].Value)
		currValue := float64(values[i].Value)
		currTime := values[i].Timestamp.Time()

		// Detect start of a spike
		if !inSpike && prevValue <= baseThreshold && currValue >= spikeThreshold {
			inSpike = true
			spikeStartTime = currTime
		}

		// Detect end of a spike
		if inSpike && currValue <= baseThreshold {
			inSpike = false
			spikeEndTime = currTime

			// We found a complete spike pattern. Let's report it.
			details := fmt.Sprintf("Instance '%s' CPU usage spiked from ~%.2f%% to %.2f%% and returned to normal.", instance, baseThreshold, spikeThreshold)
			report := models.Report{
				MetricName:      "cpu_usage",
				PatternDetected: "Sudden Morning Spike",
				Severity:        "warning",
				StartTime:       spikeStartTime,
				EndTime:         spikeEndTime,
				Details:         details,
			}

			// Save to DB
			reportID, err := database.InsertReport(report)
			if err != nil {
				log.Printf("Error saving report to DB: %v", err)
				continue
			}
			log.Printf("Detected CPU spike pattern for instance %s. Report ID: %d", instance, reportID)

			// Send notifications
			notification.SendAlert(cfg, report)
		}
	}
}