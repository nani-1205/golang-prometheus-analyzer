// internal/reporting/scheduler.go
package reporting

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/your-username/golang-prometheus-analyzer/internal/analysis"
	"github.com/your-username/golang-prometheus-analyzer/internal/config"
)

func StartScheduler(cfg *config.Config) {
	c := cron.New()

	// Schedule CPU analysis every hour at the 5-minute mark
	c.AddFunc("5 * * * *", func() {
		analysis.AnalyzeCPUUsagePattern(cfg)
	})

	// TODO: Add more analysis jobs for other metrics (memory, disk, etc.)
	// c.AddFunc("10 * * * *", func() { analysis.AnalyzeMemoryUsage(cfg) })

	log.Println("Cron scheduler initialized.")
	c.Start()
}

// NOTE: A function for a full "daily PDF report" would be added here.
// It would query the DB for all reports in the last 24 hours and use the
// PDF generator to create and email the report.