// internal/reporting/scheduler.go
package reporting

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/analysis"
	"github.com/nani-1205/golang-prometheus-analyzer/internal/config"
)

func StartScheduler(cfg *config.Config) {
	c := cron.New()

	// Schedule our analysis jobs.
	// This function will run at 5 minutes past every hour.
	c.AddFunc("5 * * * *", func() {
		// Run both types of analysis, one after the other.
		
		// 1. Run the analysis for transient spikes (low -> high -> low)
		analysis.AnalyzeCPUUsageSpikePattern(cfg)

		// 2. Run the analysis for sustained high load
		analysis.AnalyzeSustainedHighCPU(cfg)
	})

	// You can add more scheduled jobs for other metrics here.
	// For example, check memory usage at 10 minutes past the hour.
	// c.AddFunc("10 * * * *", func() { analysis.AnalyzeMemoryUsage(cfg) })

	log.Println("Cron scheduler initialized.")
	c.Start()
}