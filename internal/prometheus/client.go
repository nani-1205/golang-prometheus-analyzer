// internal/prometheus/client.go
package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/your-username/golang-prometheus-analyzer/internal/config"
)

func QueryRange(cfg *config.Config, query string, start, end time.Time, step time.Duration) (model.Value, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.PrometheusURL,
	})
	if err != nil {
		return nil, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		// For simplicity, we log warnings. In production, you might want to handle them differently.
		// log.Printf("Warnings: %v\n", warnings)
	}

	return result, nil
}