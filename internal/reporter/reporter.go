package reporter

import (
	"19_11_2026_go/internal/models"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(tasks []*models.Task) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(20)

	for _, task := range tasks {
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, fmt.Sprintf("Results for Task #%d", task.ID), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 12)
		for url, status := range task.Results {
			displayStatus := status
			if status == "available" {
				pdf.SetTextColor(34, 139, 34)
			} else {
				pdf.SetTextColor(220, 20, 60)
			}
			pdf.CellFormat(0, 8, fmt.Sprintf("  - %s: %s", url, displayStatus), "", 1, "L", false, 0, "")
		}
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(10)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
