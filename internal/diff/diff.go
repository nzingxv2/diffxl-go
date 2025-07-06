// File: internal/diff/diff.go
package diff

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/pmezard/go-difflib/difflib"
)

type DiffType int

const (
	Added DiffType = iota
	Removed
	Modified
)

type CellDiff struct {
	Row     int
	Col     int
	Type    DiffType
	Before  string
	After   string
	ColName string
}

type SheetDiffResult struct {
	SheetName string
	Diffs     []CellDiff
	AddedRows int
	RemovedRows int
}

func CompareSheet(before, after *excel.ExcelFile, sheet string, ignoreCols []string, onlyChanges bool) (SheetDiffResult, error) {
	result := SheetDiffResult{SheetName: sheet}

	// Convert sheets to CSV
	beforeCSV, err := before.ConvertSheetToCSV(sheet)
	if err != nil {
		return result, err
	}
	defer os.Remove(beforeCSV)

	afterCSV, err := after.ConvertSheetToCSV(sheet)
	if err != nil {
		return result, err
	}
	defer os.Remove(afterCSV)

	// Read CSV files
	beforeRecords, err := readCSV(beforeCSV)
	if err != nil {
		return result, err
	}

	afterRecords, err := readCSV(afterCSV)
	if err != nil {
		return result, err
	}

	// Create diff
	diff := difflib.UnifiedDiff{
		A:        recordsToLines(beforeRecords),
		B:        recordsToLines(afterRecords),
		FromFile: "Before",
		ToFile:   "After",
		Context:  0,
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return result, err
	}

	// Parse diff to structured format
	result.Diffs = parseDiff(diffText, beforeRecords[0], ignoreCols, onlyChanges)
	result.AddedRows, result.RemovedRows = calculateRowChanges(diffText)

	return result, nil
}

func readCSV(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return csv.NewReader(f).ReadAll()
}

func recordsToLines(records [][]string) []string {
	lines := make([]string, len(records))
	for i, record := range records {
		lines[i] = strings.Join(record, ",")
	}
	return lines
}

func parseDiff(diffText string, headers []string, ignoreCols []string, onlyChanges bool) []CellDiff {
	// Implementation details...
}