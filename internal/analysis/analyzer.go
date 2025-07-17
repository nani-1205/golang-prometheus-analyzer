// internal/analysis/analyzer.go
package analysis

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/common/model"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/config"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/database"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/models"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/notification"
	prom "github.com/nani-1205/golang-prometheus-analyzer/internal/prometheus"
)

// PromQL query to get average CPU usage (non-idle time) per instance
const cpuUsageQuery = `100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[2m])) * 100)`

// --- 1. ORIGINAL TRANSIENT SPIKE ANALYZER ---
// This analyzer detects a "spike and return to normal" pattern.

// AnalyzeCPUUsageSpikePattern is our example analyzer for transient spikes.
func AnalyzeCPUUsageSpikePattern(cfg *config.Config) {
	log.Println("Running Transient CPU Spike analysis...")

	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now

	result, err := prom.QueryRange(cfg, cpuUsageQuery, start, end, 1*time.Minute)
	if err != nil {
		log.Printf("Error querying Prometheus for spike analysis: %v", err)
		return
	}

	matrix, ok := result.(model.Matrix)
	if !ok {
		log.Printf("Spike Analysis Error: Expected a matrix result type, got %T", result)
		return
	}

	for _, stream := range matrix {
		instance := string(stream.Metric["instance"])
		checkAndReportSpike(cfg, stream.Values, instance)
	}
	log.Println("Transient CPU Spike analysis finished.")
}

// checkAndReportSpike is a helper to find a low -> high -> low pattern.
func checkAndReportSpike(cfg *config.Config, values []model.SamplePair, instance string) {
	const spikeThreshold = 70.0 // The peak of the spike (e.g., 70%)
	const baseThreshold = 20.0  // The normal state (e.g., 20%)

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

			details := fmt.Sprintf("Instance '%s' CPU usage spiked from ~%.2f%% to over %.2f%% and returned to normal.", instance, baseThreshold, spikeThreshold)
			report := models.Report{
				MetricName:      "cpu_transient_spike",
				PatternDetected: "Transient CPU Spike",
				Severity:        "warning",
				StartTime:       spikeStartTime,
				EndTime:         spikeEndTime,
				Details:         details,
			}

			reportID, err := database.InsertReport(report)
			if err != nil {
				log.Printf("Error saving spike report to DB: %v", err)
				continue
			}
			log.Printf("Detected CPU spike pattern for instance %s. Report ID: %d", instance, reportID)
			notification.SendAlert(cfg, report)
		}
	}
}

// --- 2. NEW SUSTAINED HIGH LOAD ANALYZER ---
// This analyzer detects when CPU usage stays high for a period. This will trigger with your load.sh script.

// AnalyzeSustainedHighCPU checks if CPU has been consistently high.
func AnalyzeSustainedHighCPU(cfg *config.Config) {
	log.Println("Running Sustained High CPU analysis...")

	// Thresholds for this specific analysis
	const highLoadThreshold = 80.0 // Trigger if CPU is above 80%
	const cooldownDuration = 30 * time.Minute // Don't re-alert for the same issue for 30 mins

	now := time.Now()
	// We only need to check the last few minutes for a current high load state
	start := now.Add(-5 * time.Minute)
	end := now

	result, err := prom.QueryRange(cfg, cpuUsageQuery, start, end, 1*time.Minute)
	if err != nil {
		log.Printf("Error querying Prometheus for high load analysis: %v", err)
		return
	}

	matrix, ok := result.(model.Matrix)
	if !ok {
		log.Printf("High Load Analysis Error: Expected a matrix result type, got %T", result)
		return
	}

	for _, stream := range matrix {
		instance := string(stream.Metric["instance"])

		// Check the most recent data point
		if len(stream.Values) == 0 {
			continue
		}
		latestSample := stream.Values[len(stream.Values)-1]
		latestValue := float64(latestSample.Value)
		latestTime := latestSample.Timestamp.Time()

		if latestValue >= highLoadThreshold {
			// ANTI-SPAM: Check if we've already reported this recently
			hasRecent, err := database.CheckForRecentReport("cpu_sustained_high", instance, cooldownDuration)
			if err != nil {
				log.Printf("Error checking for recent reports: %v", err)
				continue
			}

			if !hasRecent {
				// No recent report, so let's create a new one!
				details := fmt.Sprintf("Instance '%s' is under sustained high CPU load, currently at %.2f%% (Threshold: >%.2f%%).", instance, latestValue, highLoadThreshold)
				report := models.Report{
					MetricName:      "cpu_sustained_high",
					PatternDetected: "Sustained High CPU Load",
					Severity:        "critical",
					StartTime:       latestTime, // The event is happening now
					EndTime:         latestTime,
					Details:         details,
				}

				reportID, err := database.InsertReport(report)
				if err != nil {
					log.Printf("Error saving high load report to DB: %v", err)
					continue
				}
				log.Printf("Detected Sustained High CPU on instance %s. Report ID: %d", instance, reportID)
				notification.SendAlert(cfg, report)
			} else {
				log.Printf("Sustained High CPU on %s is ongoing, but already reported recently. Suppressing new alert.", instance)
			}
		}
	}
	log.Println("Sustained High CPU analysis finished.")
}