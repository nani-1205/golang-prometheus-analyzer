// internal/notification/notifier.go
package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/your-username/golang-prometheus-analyzer/internal/config"
	"github.com/your-username/golang-prometheus-analyzer/internal/models"
	"gopkg.in/gomail.v2"
)

// SendAlert orchestrates sending alerts to all configured channels.
func SendAlert(cfg *config.Config, report models.Report) {
	sendEmailAlert(cfg, report)
	sendGoogleChatAlert(cfg, report)
}

func sendEmailAlert(cfg *config.Config, report models.Report) {
	if cfg.SMTPUser == "" {
		return // Email not configured
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", cfg.AlertRecipientEmail)
	m.SetHeader("Subject", fmt.Sprintf("[%s] Alert: %s Detected on %s", report.Severity, report.PatternDetected, report.MetricName))
	m.SetBody("text/html", fmt.Sprintf("<h2>Analysis Report Alert</h2>"+
		"<p><strong>Metric:</strong> %s</p>"+
		"<p><strong>Pattern:</strong> %s</p>"+
		"<p><strong>Severity:</strong> %s</p>"+
		"<p><strong>Time Window:</strong> %s to %s</p>"+
		"<p><strong>Details:</strong> %s</p>",
		report.MetricName, report.PatternDetected, report.Severity,
		report.StartTime.Format(time.RFC1123), report.EndTime.Format(time.RFC1123), report.Details))

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email alert: %v", err)
	} else {
		log.Println("Email alert sent successfully.")
	}
}

func sendGoogleChatAlert(cfg *config.Config, report models.Report) {
	if cfg.GoogleChatWebhookURL == "" {
		return // Google Chat not configured
	}

	// Google Chat message format (CardsV2)
	card := map[string]interface{}{
		"cardsV2": []map[string]interface{}{{
			"cardId": "alertCard",
			"card": map[string]interface{}{
				"header": map[string]interface{}{
					"title":    fmt.Sprintf("Prometheus Analysis Alert: %s", report.PatternDetected),
					"subtitle": fmt.Sprintf("Metric: %s | Severity: %s", report.MetricName, report.Severity),
					"imageUrl": "https://cdn-icons-png.flaticon.com/512/8706/8706488.png", // A simple alert icon
				},
				"sections": []map[string]interface{}{{
					"widgets": []map[string]interface{}{
						{
							"decoratedText": map[string]interface{}{
								"topLabel": "Details",
								"text":     report.Details,
							},
						},
						{
							"decoratedText": map[string]interface{}{
								"topLabel": "Time Window",
								"text":     fmt.Sprintf("%s to %s", report.StartTime.Format(time.RFC1123), report.EndTime.Format(time.RFC1123)),
							},
						},
					},
				}},
			},
		}},
	}

	payload, err := json.Marshal(card)
	if err != nil {
		log.Printf("Error marshalling Google Chat payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", cfg.GoogleChatWebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating Google Chat request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Google Chat alert: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Google Chat webhook returned non-200 status: %s", resp.Status)
	} else {
		log.Println("Google Chat alert sent successfully.")
	}
}