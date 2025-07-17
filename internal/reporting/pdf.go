// internal/reporting/pdf.go
package reporting

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/your-username/golang-prometheus-analyzer/internal/models"
)

func GeneratePDFReport(reports []models.Report) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "Daily Analysis Report")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(20)

	// Table Header
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(30, 7, "Metric", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 7, "Pattern", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 7, "Severity", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 7, "Details", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 10)
	for _, r := range reports {
		pdf.CellFormat(30, 7, r.MetricName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 7, r.PatternDetected, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 7, r.Severity, "1", 0, "L", false, 0, "")
		pdf.MultiCell(80, 7, r.Details, "1", "L", false)
	}

	// Save file
	filename := fmt.Sprintf("report-%d.pdf", time.Now().Unix())
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}