// File: internal/output/output.go
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/nzingxv2/diffxl-go/internal/diff"
)

func GenerateTextOutput(results []diff.SheetDiffResult, outputPath string, summaryOnly bool) error {
	if outputPath == "" {
		outputPath = filepath.Join("output", fmt.Sprintf("diff_%s.txt", time.Now().Format("20060102_150405")))
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, result := range results {
		if summaryOnly {
			file.WriteString(fmt.Sprintf("Sheet: %s - Changes: %d\n", result.SheetName, len(result.Diffs)))
		} else {
			file.WriteString(fmt.Sprintf("\n===== Sheet: %s =====\n", result.SheetName))
			for _, d := range result.Diffs {
				file.WriteString(formatDiffLine(d))
			}
		}
	}

	color.Green("Diff saved to: %s", outputPath)
	return nil
}

func formatDiffLine(d diff.CellDiff) string {
	// Colorized output implementation...
}

func GenerateJSONOutput(results []diff.SheetDiffResult, outputPath string) error {
	// JSON serialization implementation...
}