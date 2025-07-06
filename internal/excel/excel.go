// File: internal/excel/excel.go
package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelFile struct {
	file     *excelize.File
	filePath string
}

func LoadExcelFile(filePath string) (*ExcelFile, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	return &ExcelFile{file: f, filePath: filePath}, nil
}

func (e *ExcelFile) Close() error {
	return e.file.Close()
}

func (e *ExcelFile) GetSheetNames() []string {
	return e.file.GetSheetList()
}

func (e *ExcelFile) ConvertSheetToCSV(sheet string) (string, error) {
	rows, err := e.file.GetRows(sheet)
	if err != nil {
		return "", fmt.Errorf("failed to get rows from sheet %s: %w", sheet, err)
	}

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s_%s_*.csv", filepath.Base(e.filePath), sheet))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	for _, row := range rows {
		cleanRow := make([]string, len(row))
		for i, cell := range row {
			cleanRow[i] = strings.ReplaceAll(cell, "\n", "\\n")
		}
		if _, err := tmpFile.WriteString(strings.Join(cleanRow, ",") + "\n"); err != nil {
			return "", fmt.Errorf("failed to write to temp file: %w", err)
		}
	}

	return tmpFile.Name(), nil
}