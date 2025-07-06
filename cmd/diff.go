// File: cmd/diff.go
package cmd

import (
	"os"
	"time"

	"github.com/nzingxv2/diffxl-go/internal/diff"
	"github.com/nzingxv2/diffxl-go/internal/excel"
	"github.com/nzingxv2/diffxl-go/internal/output"
	"github.com/nzingxv2/diffxl-go/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	beforeFile  string
	afterFile   string
	onlyChanges bool
	ignoreCols  []string
	summary     bool
	sheetName   string
	outputPath  string
	jsonOutput  bool
)

var diffCmd = &cobra.Command{
	Use:   "diff [flags]",
	Short: "Compare two Excel files",
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		
		// Validate input
		if beforeFile == "" || afterFile == "" {
			color.Red("Both --before and --after are required")
			os.Exit(1)
		}

		// Load Excel files
		before, err := excel.LoadExcelFile(beforeFile)
		utils.HandleError(err, "Failed to load before file")
		defer before.Close()

		after, err := excel.LoadExcelFile(afterFile)
		utils.HandleError(err, "Failed to load after file")
		defer after.Close()

		// Get sheets to compare
		sheets := utils.GetSheetsToCompare(before.GetSheetNames(), after.GetSheetNames(), sheetName)
		if len(sheets) == 0 {
			color.Yellow("No sheets to compare")
			return
		}

		// Process diff for each sheet
		var diffResults []diff.SheetDiffResult
		for _, sheet := range sheets {
			result, err := diff.CompareSheet(before, after, sheet, ignoreCols, onlyChanges)
			utils.HandleError(err, "Diff failed for sheet: "+sheet)
			diffResults = append(diffResults, result)
		}

		// Generate output
		if jsonOutput {
			err = output.GenerateJSONOutput(diffResults, outputPath)
		} else {
			err = output.GenerateTextOutput(diffResults, outputPath, summary)
		}
		utils.HandleError(err, "Failed to generate output")

		// Show summary
		if summary {
			stats := diff.CalculateSummaryStats(diffResults)
			color.Cyan("\nSummary:")
			color.Cyan("  Sheets compared: %d", stats.SheetsCompared)
			color.Cyan("  Changed cells: %d", stats.ChangedCells)
			color.Cyan("  Added rows: %d", stats.AddedRows)
			color.Cyan("  Removed rows: %d", stats.RemovedRows)
			color.Cyan("  Execution time: %v", time.Since(startTime).Round(time.Millisecond))
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&beforeFile, "before", "b", "", "Path to the 'before' Excel file (required)")
	diffCmd.Flags().StringVarP(&afterFile, "after", "a", "", "Path to the 'after' Excel file (required)")
	diffCmd.Flags().BoolVar(&onlyChanges, "only-changes", false, "Show only changed cells")
	diffCmd.Flags().StringSliceVar(&ignoreCols, "ignore-columns", []string{}, "Columns to ignore (comma separated)")
	diffCmd.Flags().BoolVar(&summary, "summary", false, "Show summary statistics")
	diffCmd.Flags().StringVar(&sheetName, "sheet", "", "Specific sheet to compare (default: all sheets)")
	diffCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: output/diff_<timestamp>.txt)")
	diffCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	diffCmd.MarkFlagRequired("before")
	diffCmd.MarkFlagRequired("after")
}